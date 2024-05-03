import argparse
import io
import json
import logging
import threading
import uuid
from concurrent import futures
from concurrent.futures import ThreadPoolExecutor, wait
from typing import Iterable

import grpc
import torch
from google.protobuf import empty_pb2
from torch.nn import Module

import client_pb2
import client_pb2_grpc
import config
import edge_pb2_grpc
from consulx import Consul
from mq import RabbitMQConnection
from util import get_free_port

conf = config.load_configuration("config.yaml")
client_credentials = config.load_client_credentials("../../internal/authorization")
server_credentials = config.load_server_credentials('../../internal/authorization')
service_conf = conf['service']
consul_conf = conf['consul']
rabbitmq_conf = conf['rabbitmq']
client_conf = conf['client']

# Set up basic logging
logging.basicConfig(format='%(asctime)s - %(levelname)s - %(message)s', level=logging.INFO, datefmt='%Y-%m-%d %H:%M:%S')
logging.getLogger("pika").setLevel(logging.WARNING)

STREAM_PART_SIZE = 2 * 1024 * 1024
MAX_MESSAGE_LENGTH = 1024 * 1024 * 1024
options = [('grpc.max_send_message_length', MAX_MESSAGE_LENGTH),
           ('grpc.max_receive_message_length', MAX_MESSAGE_LENGTH), ('grpc.enable_retries', 1)]
COMMON_NAME = "ccc"
consul = Consul(consul_host=consul_conf['host'], consul_port=consul_conf['port'])
rabbitmq = RabbitMQConnection(host=rabbitmq_conf['host'], port=rabbitmq_conf['port'])
rabbitmq_connection = rabbitmq.open_connection()


def consume_message(ch, method, properties, body):
    try:
        data = json.loads(body)
        services = data['services']
        configuration = data['conf']
        # 自定义模型参数
        with open('test.pt', 'rb') as f:
            new_pt = f.read()
        for _ in range(configuration['global_epochs']):
            with ThreadPoolExecutor(max_workers=len(services)) as executor:
                tasks = [
                    executor.submit(train_model, address=service, pt=new_pt, configuration=json.dumps(configuration))
                    for
                    service in services]
                wait(tasks)
                models = [task.result() for task in tasks if task.result() is not None]
                new_pt = aggregate(configuration, models)
        ch.basic_ack(delivery_tag=method.delivery_tag)
    except Exception as ex:
        logging.error(f"Failed to process message: {ex}")
        # 处理失败，拒绝消息并决定是否重新投递
        ch.basic_nack(delivery_tag=method.delivery_tag, requeue=True)
    finally:
        executor.shutdown(wait=True)


def send_stream_data(pt: bytes, configuration: str):
    """Iterators send large files"""
    file_size = len(pt)
    with io.BytesIO(pt) as stream:
        while True:
            chunk = stream.read(STREAM_PART_SIZE)
            if chunk:
                yield client_pb2.TrainModelRequest(pt=chunk, conf=configuration, size=file_size)
            else:
                break


def train_model(address: str, pt: bytes, configuration: str) -> bytes:
    try:
        channel_options = options + [('grpc.ssl_target_name_override', COMMON_NAME)]
        with grpc.intercept_channel(
                grpc.secure_channel(address, client_credentials, options=channel_options)) as channel:
            stub = client_pb2_grpc.ClientServiceStub(channel)
            request = send_stream_data(pt=pt, configuration=configuration)
            data = bytearray()
            for response in stub.TrainModel(request):
                size = response.size
                data.extend(response.pt)
                if len(data) == size:
                    break
            return torch.load(io.BytesIO(data))
    except Exception as exception:
        logging.error(f"Failed to train model with {address}: {exception}")
        raise


def aggregate(configuration: dict, modules: Iterable[Module]) -> bytes:
    # Placeholder for aggregation logic
    # This needs to be implemented based on your specific aggregation algorithm
    logging.info("Aggregating models...")
    buf = io.BytesIO()
    torch.save(modules, buf)
    return buf.getvalue()


class EdgeServicer(edge_pb2_grpc.EdgeServiceServicer):

    def TrainTask(self, request, context):
        try:
            conf_json = json.loads(request.conf)
            client_num = conf_json["clients"]
            services = consul.get_services(client_conf['name'], client_num)
            if len(services) < client_num:
                context.set_details('Not enough available clients.')
                context.set_code(grpc.StatusCode.UNAVAILABLE)
                logging.error('Not enough available clients.')
                return empty_pb2.Empty()

            rabbitmq.send_message_to_queue(rabbitmq_connection, services, conf_json)
            return empty_pb2.Empty()
        except Exception as ex:
            logging.error(f"Error during training: {ex}")
            context.set_details(f'An error occurred during model training: {ex}')
            context.set_code(grpc.StatusCode.INTERNAL)
            return empty_pb2.Empty()


if __name__ == '__main__':
    parser = argparse.ArgumentParser(description="Start the server with specified port")
    parser.add_argument('-p', '--port', type=int, help='The port to start the server on')
    args = parser.parse_args()
    port = args.port if args.port else get_free_port()
    server = None
    service_id = str(uuid.uuid4())
    try:
        server = grpc.server(futures.ThreadPoolExecutor(max_workers=1000), options=options)
        edge_pb2_grpc.add_EdgeServiceServicer_to_server(EdgeServicer(), server)
        server_address = f"{service_conf['host']}:{port}"
        server.add_secure_port(server_address, server_credentials)
        server.start()
        logging.info(f"Service started at {server_address}")
        consul.register(address=service_conf['host'], port=port, service_name=service_conf['name'],
                        tags=service_conf['tags'], service_id=service_id)
        logging.info(f"Service registered with Consul successfully.")
        for _ in range(5):
            rabbitmq_consumer_thread = threading.Thread(target=rabbitmq.consume_messages_from_queue,
                                                        args=(consume_message,))
            rabbitmq_consumer_thread.daemon = True
            rabbitmq_consumer_thread.start()
        server.wait_for_termination()
    except KeyboardInterrupt:
        server.stop(0)
        logging.info("Server stopped by KeyboardInterrupt.")
        consul.deregister(server_id=service_id)
        logging.info(f"Service unregistered from Consul successfully.")
        rabbitmq.close_connection(rabbitmq_connection)
    except Exception as e:
        logging.error(f"An error occurred: {e}")
        raise
