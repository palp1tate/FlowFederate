import argparse
import io
import json
import logging
import uuid
from concurrent import futures
from datetime import datetime

import grpc
import torch
from sqlalchemy.orm import sessionmaker
from torch.utils.data import DataLoader

import client_pb2
import client_pb2_grpc
from consulx import Consul
from datasets import get_dataset
from get_model import get_model
from internal.model.table import Client
from util import *

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
server = None
service_id = str(uuid.uuid4())
consul = Consul(consul_host=consul_conf["host"], consul_port=consul_conf["port"])
engine = init_engine()


def train_model(pt: bytes, configuration: dict, task_id: int) -> bytes:
    try:
        model_name = configuration["model_name"]
        dataset_type = configuration["type"]
        local_epochs = configuration["local_epochs"]
        batch_size = configuration["batch_size"]
        learning_rate = configuration["lr"]
        momentum = configuration["momentum"]
        new_session = sessionmaker(engine)
        with new_session() as session:
            client = (
                session.query(Client)
                .filter_by(task_id=task_id, client_id=service_id)
                .first()
            )
            if not client:
                client = Client(
                    task_id=task_id,
                    client_id=service_id,
                    created_at=datetime.now(),
                    updated_at=datetime.now(),
                    model=model_name,
                    dataset=dataset_type,
                    type="通用模型",
                    status="正常",
                    current_round=0,
                    total_round=local_epochs,
                    accuracy=None,
                    loss=None,
                    cpu=get_cpu_usage(),
                    memory=get_memory_usage(),
                    disk=get_disk_usage(),
                    progress="0%",
                )
                session.add(client)
                session.commit()

        logging.info("Training model...")
        state_dict = torch.load(io.BytesIO(pt))

        logging.info(
            f"Model: {model_name}, Dataset: {dataset_type}, Local epochs: {local_epochs}, Batch size: {batch_size}, "
            f"Learning rate: {learning_rate}, Momentum: {momentum}"
        )

        local_model = get_model(model_name)

        local_model.load_state_dict(state_dict)
        device = "cuda" if torch.cuda.is_available() else "cpu"
        local_model = local_model.to(device)

        train_datasets, eval_datasets = get_dataset("./data/", dataset_type)
        train_loader = DataLoader(train_datasets, batch_size=batch_size, shuffle=True)
        eval_loader = DataLoader(eval_datasets, batch_size=batch_size, shuffle=False)

        optimizer_name = configuration.get("optimizer", "sgd")  # 默认使用SGD
        if optimizer_name == "sgd":
            optimizer = torch.optim.SGD(
                local_model.parameters(), lr=learning_rate, momentum=momentum
            )
        elif optimizer_name == "adam":
            optimizer = torch.optim.Adam(
                local_model.parameters(), lr=learning_rate, betas=(0.9, 0.999)
            )
        elif optimizer_name == "rmsprop":
            optimizer = torch.optim.RMSprop(
                local_model.parameters(), lr=learning_rate, alpha=0.99
            )
        else:
            raise ValueError(f"Unsupported optimizer type: {optimizer_name}")
        logging.info(f"Optimizer: {optimizer_name}")

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
        logging.info(f"Loss function: {loss_function_name}")

        for epoch in range(local_epochs):
            local_model.train()
            total_train_loss = 0
            for data, target in train_loader:
                data, target = data.to(device), target.to(device)
                optimizer.zero_grad()
                _, output = local_model(data)
                loss = criterion(output, target)
                loss.backward()
                optimizer.step()
                total_train_loss += loss.item()

            # 计算平均训练损失
            avg_train_loss = total_train_loss / len(train_loader)

            # 评估模型
            local_model.eval()
            total_eval_loss = 0
            correct = 0
            with torch.no_grad():
                for data, target in eval_loader:
                    data, target = data.to(device), target.to(device)
                    _, output = local_model(data)
                    loss = criterion(output, target)
                    total_eval_loss += loss.item()
                    pred = output.argmax(dim=1, keepdim=True)  # 获取预测结果
                    correct += pred.eq(target.view_as(pred)).sum().item()

            # 计算平均验证损失和精度
            avg_eval_loss = total_eval_loss / len(eval_loader)
            accuracy = correct / len(eval_loader.dataset)

            print(
                f"Epoch {epoch + 1}/{local_epochs}, Train Loss: {avg_train_loss:.4f}, Eval Loss: {avg_eval_loss:.4f}, "
                f"Accuracy: {accuracy:.4f}"
            )

            with new_session() as session:
                client = (
                    session.query(Client)
                    .filter_by(task_id=task_id, client_id=service_id)
                    .first()
                )
                if client:
                    client.updated_at = datetime.now()
                    client.current_round = epoch + 1
                    client.accuracy = accuracy
                    client.loss = avg_eval_loss
                    client.cpu = get_cpu_usage()
                    client.memory = get_memory_usage()
                    client.disk = get_disk_usage()
                    client.progress = f"{(epoch + 1) * 100 // local_epochs}%"
                    client.status = "正常"
                    session.commit()

        trained_model = local_model.state_dict()
        buf = io.BytesIO()
        torch.save(trained_model, buf)
        return buf.getvalue()
    except Exception as exc:
        new_session = sessionmaker(engine)
        with new_session() as session:
            client = (
                session.query(Client)
                .filter_by(task_id=task_id, client_id=service_id)
                .first()
            )
            if client is not None:
                client.status = "异常"
                client.updated_at = datetime.now()
                session.commit()
        logging.error(f"Model training failed: {exc}")
        raise


class ClientServicer(client_pb2_grpc.ClientServiceServicer):

    def TrainModel(self, request_iterator, context):
        try:
            pt = bytearray()
            conf_json = None
            task_id = None
            for TrainRequest in request_iterator:
                pt_size = TrainRequest.size
                pt_data = TrainRequest.pt
                if conf_json is None:
                    conf_json = json.loads(TrainRequest.conf)
                if task_id is None:
                    task_id = TrainRequest.task_id
                pt.extend(pt_data)
                if len(pt) == pt_size:
                    break
            logging.info(f"Preparing to train model...")
            new_pt = train_model(pt=pt, configuration=conf_json, task_id=task_id)
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
