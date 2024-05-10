import argparse
import io
import json
import threading
import uuid
from concurrent import futures
from concurrent.futures import ThreadPoolExecutor, wait
from datetime import datetime

import grpc
from google.protobuf import empty_pb2
from sqlalchemy.orm import sessionmaker
from torch.utils.data import DataLoader

import client_pb2
import client_pb2_grpc
import edge_pb2_grpc
from aggregate import *
from consulx import Consul
from datasets import get_dataset
from get_model import get_model
from internal.model.table import Server
from internal.model.table import Task
from mq import RabbitMQConnection
from util import *

conf = config.load_configuration("service/edge/config.yaml")
client_credentials = config.load_client_credentials("internal/authorization")
server_credentials = config.load_server_credentials("internal/authorization")
service_conf = conf["service"]
consul_conf = conf["consul"]
rabbitmq_conf = conf["rabbitmq"]
client_conf = conf["client"]

# Set up basic logging
logging.basicConfig(
    format="%(asctime)s - %(levelname)s - %(message)s",
    level=logging.INFO,
    datefmt="%Y-%m-%d %H:%M:%S",
)
logging.getLogger("pika").setLevel(logging.WARNING)

STREAM_PART_SIZE = 2 * 1024 * 1024
MAX_MESSAGE_LENGTH = 1024 * 1024 * 1024
options = [
    ("grpc.max_send_message_length", MAX_MESSAGE_LENGTH),
    ("grpc.max_receive_message_length", MAX_MESSAGE_LENGTH),
    ("grpc.enable_retries", 1),
]
COMMON_NAME = "ccc"
server = None
service_id = str(uuid.uuid4())
consul = Consul(consul_host=consul_conf["host"], consul_port=consul_conf["port"])
rabbitmq = RabbitMQConnection(host=rabbitmq_conf["host"], port=rabbitmq_conf["port"])
rabbitmq_connection = rabbitmq.open_connection()
engine = init_engine()
device = "cuda" if torch.cuda.is_available() else "cpu"


def consume_message(ch, method, properties, body):
    try:
        data = json.loads(body)
        services = data["services"]
        configuration = data["conf"]
        user_name = data["user_name"]

        # 新建训练任务
        global_epochs = configuration["global_epochs"]
        local_epochs = configuration["local_epochs"]
        model_name = configuration["model_name"]
        dataset_type = configuration["type"]

        new_session = sessionmaker(engine)
        with new_session() as session:
            created_at = datetime.now()
            updated_at = datetime.now()
            task = Task(
                user_name=user_name,
                created_at=created_at,
                updated_at=updated_at,
                model=model_name,
                dataset=dataset_type,
                type="通用模型",
                status="正常",
                progress="0%",
            )
            session.add(task)
            session.commit()
            task_id = task.id
            edge_server = (
                session.query(Server)
                .filter_by(task_id=task_id, server_id=service_id)
                .first()
            )
            if not edge_server:
                edge_server = Server(
                    task_id=task_id,
                    server_id=service_id,
                    created_at=created_at,
                    updated_at=updated_at,
                    model=model_name,
                    dataset=dataset_type,
                    type="通用模型",
                    status="正常",
                    current_round=0,
                    total_round=global_epochs,
                    cpu=get_cpu_usage(),
                    memory=get_memory_usage(),
                    disk=get_disk_usage(),
                    progress="0%",
                )
                session.add(edge_server)
            session.commit()
        print(f"Task {task_id} created successfully.")

        # 自定义模型参数
        local_model = get_model(model_name)
        local_model_state_dict = local_model.state_dict()

        # Save state_dict to a bytes buffer
        buffer = io.BytesIO()
        torch.save(local_model_state_dict, buffer)

        new_pt = buffer.getvalue()
        logging.info(f"Training model for {local_epochs} epochs...")
        for epoch in range(global_epochs):
            logging.info(f"Epoch {epoch + 1}/{global_epochs}")
            with ThreadPoolExecutor(max_workers=len(services)) as executor:
                tasks = [
                    executor.submit(
                        train_model,
                        address=service,
                        pt=new_pt,
                        configuration=json.dumps(configuration),
                        task_id=task_id,
                    )
                    for service in services
                ]
                wait(tasks)
                models = [
                    torch.load(io.BytesIO(task.result()))
                    for task in tasks
                    if task.result() is not None
                ]
                modules = []
                for state_dict in models:
                    model = get_model(model_name)
                    model.load_state_dict(state_dict)
                    modules.append(model)
                logging.info("Aggregating models...")
                new_pt, acc, loss = aggregate(configuration, modules, new_pt)
                progress = f"{(epoch + 1) * 100 // global_epochs}%"  # 计算进度
                with new_session() as session:
                    # 更新任务状态
                    updated_at = datetime.now()
                    task = session.get(Task, task_id)
                    if task:
                        task.accuracy = acc
                        task.loss = loss
                        task.progress = progress
                        task.updated_at = updated_at
                        task.status = "正常"

                    # 更新边缘服务器状态
                    edge_server = (
                        session.query(Server)
                        .filter_by(task_id=task_id, server_id=service_id)
                        .first()
                    )
                    if edge_server:
                        edge_server.updated_at = updated_at
                        edge_server.current_round = epoch + 1
                        edge_server.accuracy = acc
                        edge_server.loss = loss
                        edge_server.cpu = get_cpu_usage()
                        edge_server.memory = get_memory_usage()
                        edge_server.disk = get_disk_usage()
                        edge_server.progress = progress
                        edge_server.status = "正常"
                    session.commit()

        ch.basic_ack(delivery_tag=method.delivery_tag)
    except Exception as ex:
        logging.error(f"Failed to process message: {ex}")
        new_session = sessionmaker(engine)
        with new_session() as session:
            task = session.get(Task, task_id)
            if task:
                task.status = "异常"
                task.updated_at = datetime.now()
                session.commit()
            edge_server = (
                session.query(Server)
                .filter_by(task_id=task_id, server_id=service_id)
                .first()
            )
            if edge_server:
                edge_server.status = "异常"
                edge_server.updated_at = datetime.now()
                session.commit()
            session.commit()
        # 处理失败，拒绝消息并决定是否重新投递
        ch.basic_nack(delivery_tag=method.delivery_tag, requeue=True)
    finally:
        executor.shutdown(wait=True)


def send_stream_data(pt: bytes, configuration: str, task_id: int):
    """Iterators send large files"""
    file_size = len(pt)
    with io.BytesIO(pt) as stream:
        while True:
            chunk = stream.read(STREAM_PART_SIZE)
            if chunk:
                yield client_pb2.TrainModelRequest(
                    pt=chunk, conf=configuration, size=file_size, task_id=task_id
                )
            else:
                break


def train_model(address: str, pt: bytes, configuration: str, task_id: int) -> bytes:
    try:
        channel_options = options + [
            ("grpc.ssl_target_name_override", COMMON_NAME),
            ("grpc.keepalive_time_ms", 0),
        ]
        with grpc.intercept_channel(
            grpc.secure_channel(address, client_credentials, options=channel_options)
        ) as channel:
            stub = client_pb2_grpc.ClientServiceStub(channel)
            request = send_stream_data(
                pt=pt, configuration=configuration, task_id=task_id
            )
            data = bytearray()
            for response in stub.TrainModel(request):
                size = response.size
                data.extend(response.pt)
                if len(data) == size:
                    break
            return data
    except Exception as exception:
        logging.error(f"Failed to train model with {address}: {exception}")
        raise


def aggregate(
    configuration: dict, modules: Iterable[torch.nn.Module], pt: bytes
) -> (bytes, float, float):
    logging.info("Aggregate Function Aggregating models...")

    # 加载全局模型的当前状态
    state_dict = torch.load(io.BytesIO(pt))
    global_model = get_model(configuration["model_name"])
    ori_state_dict = global_model.state_dict()
    for key in list(state_dict.keys()):
        if key not in ori_state_dict:
            del state_dict[key]
    global_model.load_state_dict(state_dict)
    logging.info("Global model loaded successfully.")
    global_model.to(device)

    batch_size = configuration["batch_size"]
    dataset_type = configuration["type"]
    _, eval_datasets = get_dataset("./data/", dataset_type)
    logging.info("Dataset loaded successfully.")
    eval_loader = DataLoader(eval_datasets, batch_size=batch_size, shuffle=False)
    logging.info("Dataloader loaded successfully.")

    # 确定聚合方法
    logging.info(f"Aggregating using {configuration.get('method')} method...")
    if configuration.get("method") == "krum":
        selected_model = krum_aggregate(modules, len(modules), device)
        global_model.load_state_dict(selected_model)
    elif configuration.get("method") == "fedavg":
        avg_state_dict = fedavg_aggregate(modules, device)
        global_model.load_state_dict(avg_state_dict)
    elif configuration.get("method") == "median":
        median_params = median_aggregate(modules, device)
        global_model.load_state_dict(median_params)
    elif configuration.get("method") == "pefl":
        weighted_params = pefl_aggregate(modules, device)
        global_model.load_state_dict(weighted_params)
    elif configuration.get("method") == "trimmed_mean":
        poisoner_nums = configuration.get("poisoner_nums", 0)
        candidates = len(modules)
        logging.info(f"Poisoner nums: {poisoner_nums}, Candidates: {candidates}")
        if candidates - 2 * poisoner_nums > 0:
            trimmed_mean_params = trimmed_mean_aggregate(modules, device, poisoner_nums)
            global_model.load_state_dict(trimmed_mean_params)
        else:
            logging.error(
                "Not enough candidates after trimming to perform aggregation."
            )
    elif configuration.get("method") == "shieldfl":
        aggregated_state_dict = shieldfl_aggregate(modules, global_model, device)
        global_model.load_state_dict(aggregated_state_dict)
    else:
        avg_state_dict = fedavg_aggregate(modules, device)
        global_model.load_state_dict(avg_state_dict)
    logging.info("Global model aggregated successfully.")

    # 评估更新后的全局模型
    global_model.eval()
    logging.info("Getting eval loss and accuracy...")
    loss_function_name = configuration.get(
        "loss_function", "cross_entropy"
    )  # 默认为交叉熵损失
    if loss_function_name == "cross_entropy":
        criterion = torch.nn.CrossEntropyLoss()
    elif loss_function_name == "mse":
        criterion = torch.nn.MSELoss()
    elif loss_function_name == "nll":
        criterion = torch.nn.NLLLoss()
    else:
        raise ValueError(f"Unsupported loss function type: {loss_function_name}")
    total_loss = 0
    correct = 0
    total = 0
    with torch.no_grad():
        for data, target in eval_loader:
            data, target = data.to(device), target.to(device)
            _, outputs = global_model(data)
            loss = criterion(outputs, target)
            total_loss += loss.item()
            _, predicted = torch.max(outputs.data, 1)
            total += target.size(0)
            correct += (predicted == target).sum().item()

    avg_loss = total_loss / len(eval_loader)
    accuracy = correct / total
    print(f"Global Eval Loss: {avg_loss:.4f}, Accuracy: {accuracy:.4f}")

    # 保存并返回更新后的全局模型参数
    buf = io.BytesIO()
    torch.save(global_model.state_dict(), buf)
    return buf.getvalue(), accuracy, avg_loss


class EdgeServicer(edge_pb2_grpc.EdgeServiceServicer):

    def TrainTask(self, request, context):
        try:
            conf_json = json.loads(request.conf)
            user_name = request.user_name
            client_num = conf_json["clients"]
            services = consul.get_services(client_conf["name"], client_num)
            if len(services) < client_num:
                context.set_details("Not enough available clients.")
                context.set_code(grpc.StatusCode.UNAVAILABLE)
                logging.error("Not enough available clients.")
                return empty_pb2.Empty()

            rabbitmq.send_message_to_queue(
                rabbitmq_connection, services, conf_json, user_name
            )
            return empty_pb2.Empty()
        except Exception as ex:
            logging.error(f"Error during training: {ex}")
            context.set_details(f"An error occurred during model training: {ex}")
            context.set_code(grpc.StatusCode.INTERNAL)
            return empty_pb2.Empty()


if __name__ == "__main__":
    parser = argparse.ArgumentParser(description="Start the server with specified port")
    parser.add_argument(
        "-p", "--port", type=int, help="The port to start the server on"
    )
    args = parser.parse_args()
    port = args.port if args.port else get_free_port()
    try:
        server = grpc.server(
            futures.ThreadPoolExecutor(max_workers=1000), options=options
        )
        edge_pb2_grpc.add_EdgeServiceServicer_to_server(EdgeServicer(), server)
        server_address = f"[::]:{port}"
        server.add_secure_port(server_address, server_credentials)
        server_host = get_ip_address()
        server.start()
        logging.info(f"Service started at {server_host}:{port}")
        consul.register(
            address=server_host,
            port=port,
            service_name=service_conf["name"],
            tags=service_conf["tags"],
            service_id=service_id,
        )
        logging.info(f"Service registered with Consul successfully.")
        for _ in range(5):
            rabbitmq_consumer_thread = threading.Thread(
                target=rabbitmq.consume_messages_from_queue, args=(consume_message,)
            )
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
