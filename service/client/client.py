import argparse
import io
import json
import logging
import uuid
from concurrent import futures

import grpc
import torch

import client_pb2
import client_pb2_grpc
import config
from consulx import Consul
from util import get_free_port

conf = config.load_configuration("config.yaml")
client_credentials = config.load_client_credentials("../../internal/authorization")
server_credentials = config.load_server_credentials("../../internal/authorization")
service_conf = conf["service"]
consul_conf = conf["consul"]

# Set up basic logging
logging.basicConfig(
    format="%(asctime)s - %(levelname)s - %(message)s",
    level=logging.INFO,
    datefmt="%Y-%m-%d %H:%M:%S",
)

STREAM_PART_SIZE = 2 * 1024 * 1024
MAX_MESSAGE_LENGTH = 1024 * 1024 * 1024
options = [
    ("grpc.max_send_message_length", MAX_MESSAGE_LENGTH),
    ("grpc.max_receive_message_length", MAX_MESSAGE_LENGTH),
    ("grpc.enable_retries", 1),
]
COMMON_NAME = "ccc"
consul = Consul(consul_host=consul_conf["host"], consul_port=consul_conf["port"])


def train_model(pt: bytes, configuration: dict) -> bytes:
    """Simulates model training. This function should be fully implemented based on actual training logic."""
    try:
        logging.info("Training model...")
        model = torch.load(io.BytesIO(pt))
        # Placeholder for actual training logic
        trained_model = model  # This should be replaced with actual training logic
        buf = io.BytesIO()
        torch.save(trained_model, buf)
        return buf.getvalue()
    except Exception as exc:
        logging.error(f"Model training failed: {exc}")
        raise


class ClientServicer(client_pb2_grpc.ClientServiceServicer):
    def TrainModel(self, request_iterator, context):
        try:
            pt = bytearray()
            conf_json = None
            for TrainRequest in request_iterator:
                pt_size = TrainRequest.size
                pt_data = TrainRequest.pt
                if conf_json is None:
                    conf_json = json.loads(TrainRequest.conf)
                pt.extend(pt_data)
                if len(pt) == pt_size:
                    break
            new_pt = train_model(pt=pt, configuration=conf_json)
            file_size = len(new_pt)  # 获取pt的大小
            with io.BytesIO(new_pt) as stream:
                while True:
                    chunk = stream.read(STREAM_PART_SIZE)
                    if chunk:
                        yield client_pb2.TrainModelResponse(pt=chunk, size=file_size)
                    else:
                        break
        except json.JSONDecodeError as exception:
            logging.error(f"JSON decoding error: {exception}")
            context.set_details("Invalid JSON configuration.")
            context.set_code(grpc.StatusCode.INVALID_ARGUMENT)
            return client_pb2.TrainModelResponse()
        except Exception as exception:
            logging.error(f"Error during model training: {exception}")
            context.set_details(f"An error occurred during model training: {exception}")
            context.set_code(grpc.StatusCode.INTERNAL)
            return client_pb2.TrainModelResponse()


if __name__ == "__main__":
    parser = argparse.ArgumentParser(description="Start the server with specified port")
    parser.add_argument(
        "-p", "--port", type=int, help="The port to start the server on"
    )
    args = parser.parse_args()
    port = args.port if args.port else get_free_port()
    server = None
    service_id = str(uuid.uuid4())
    try:
        server = grpc.server(
            futures.ThreadPoolExecutor(max_workers=1000), options=options
        )
        client_pb2_grpc.add_ClientServiceServicer_to_server(ClientServicer(), server)
        server_address = f"{service_conf['host']}:{port}"
        server.add_secure_port(server_address, server_credentials)
        server.start()
        logging.info(f"Service started at {server_address}")
        consul.register(
            address=service_conf["host"],
            port=port,
            service_name=service_conf["name"],
            tags=service_conf["tags"],
            service_id=service_id,
        )
        logging.info(f"Service registered with Consul successfully.")
        server.wait_for_termination()
    except KeyboardInterrupt:
        server.stop(0)
        logging.info("Server stopped by KeyboardInterrupt.")
        consul.deregister(server_id=service_id)
        logging.info(f"Service unregistered from Consul successfully.")
    except Exception as e:
        logging.error(f"An error occurred: {e}")
        raise
