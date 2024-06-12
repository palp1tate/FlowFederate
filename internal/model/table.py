from sqlalchemy import (
    Column,
    String,
    BigInteger,
    Float,
    DateTime,
    ForeignKey,
    Integer,
    CheckConstraint,
)
from sqlalchemy.ext.declarative import declarative_base
from sqlalchemy.orm import relationship

Base = declarative_base()


class UserInfo(Base):
    __tablename__ = "user_info"
    uuid = Column(String(255), nullable=False)
    user_name = Column(String(255), primary_key=True)
    password = Column(String(255), nullable=False)
    role = Column(Integer, nullable=False, default=1)
    state = Column(Integer, nullable=False, default=0)
    create_time = Column(String(255), nullable=False)

    __table_args__ = (
        CheckConstraint("role IN (0, 1)", name="user_info_chk_1"),
        CheckConstraint("state IN (0, 1)", name="user_info_chk_2"),
    )


class Task(Base):
    __tablename__ = "task"
    id = Column(BigInteger, primary_key=True, autoincrement=True)
    user_name = Column(String(256))
    created_at = Column(DateTime)
    updated_at = Column(DateTime)
    model = Column(String(256))
    dataset = Column(String(256))
    type = Column(String(256))
    status = Column(String(256))
    progress = Column(String(256))
    accuracy = Column(String(2048))
    loss = Column(String(2048))

    clients = relationship(
        "Client", back_populates="task", cascade="all, delete-orphan"
    )
    servers = relationship(
        "Server", back_populates="task", cascade="all, delete-orphan"
    )


class Server(Base):
    __tablename__ = "server"
    id = Column(BigInteger, primary_key=True, autoincrement=True)
    server_id = Column(String(256))
    task_id = Column(
        BigInteger, ForeignKey("task.id", ondelete="CASCADE", onupdate="CASCADE")
    )
    created_at = Column(DateTime)
    updated_at = Column(DateTime)
    model = Column(String(256))
    dataset = Column(String(256))
    type = Column(String(256))
    status = Column(String(256))
    current_round = Column(BigInteger)
    total_round = Column(BigInteger)
    progress = Column(String(256))
    accuracy = Column(Float)
    loss = Column(Float)
    cpu = Column(String(256))
    memory = Column(String(256))
    disk = Column(String(256))

    task = relationship("Task", back_populates="servers")


class Client(Base):
    __tablename__ = "client"
    id = Column(BigInteger, primary_key=True, autoincrement=True)
    client_id = Column(String(256))
    task_id = Column(
        BigInteger, ForeignKey("task.id", ondelete="CASCADE", onupdate="CASCADE")
    )
    created_at = Column(DateTime)
    updated_at = Column(DateTime)
    model = Column(String(256))
    dataset = Column(String(256))
    type = Column(String(256))
    status = Column(String(256))
    current_round = Column(BigInteger)
    total_round = Column(BigInteger)
    progress = Column(String(256))
    accuracy = Column(Float)
    loss = Column(Float)
    cpu = Column(String(256))
    memory = Column(String(256))
    disk = Column(String(256))

    task = relationship("Task", back_populates="clients")
