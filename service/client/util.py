import os
import socket

import psutil
from sqlalchemy import create_engine

import config


def get_free_port():
    with socket.socket(socket.AF_INET, socket.SOCK_STREAM) as s:
        s.bind(("", 0))
        return s.getsockname()[1]


def get_ip_address():
    # 检测操作系统
    if os.name == "nt":  # Windows
        # 获取本机电脑名
        hostname = socket.gethostname()
        # 获取本机ip
        ip_address = socket.gethostbyname(hostname)
        return ip_address
    elif os.name == "posix":  # Linux 或 macOS
        try:
            with socket.socket(socket.AF_INET, socket.SOCK_DGRAM) as s:
                s.connect(("8.8.8.8", 80))
                ip_address = s.getsockname()[0]
                return ip_address
        except Exception:
            return "127.0.0.1"
    else:
        return "Unsupported OS"


def get_cpu_usage() -> str:
    cpu_percent = psutil.cpu_percent(interval=1, percpu=False)
    return f"{cpu_percent}%"


def get_memory_usage() -> str:
    memory = psutil.virtual_memory()
    return f"{memory.percent}%"


def get_disk_usage() -> str:
    disk_usage = psutil.disk_usage("/")
    return f"{disk_usage.percent}%"


def init_engine():
    conf = config.load_configuration("service/client/config.yaml")
    mysql_conf = conf["mysql"]

    DATABASE_URI = (
        f"mysql+mysqlconnector://{mysql_conf['user']}:{mysql_conf['password']}@"
        f"{mysql_conf['host']:{mysql_conf['port']}}/{mysql_conf['database']}"
    )

    engine = create_engine(DATABASE_URI, pool_recycle=3600, future=True)
    return engine
