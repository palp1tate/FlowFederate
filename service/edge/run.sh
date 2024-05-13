#!/bin/sh

# 尝试连接到 MySQL
while ! nc -z mysql 3306; do
  echo "Waiting for MySQL..."
  sleep 1
done

echo "MySQL is up..."

# 尝试连接到 Consul
while ! nc -z consul 8500; do
  echo "Waiting for Consul..."
  sleep 1
done

echo "Consul is up..."

# 尝试连接到 RabbitMQ
while ! nc -z rabbitmq 5672; do
  echo "Waiting for RabbitMQ..."
  sleep 1
done

echo "RabbitMQ is up..."

# 设置 PYTHONPATH 环境变量
export PYTHONPATH=./

# 执行 python 程序
python ./service/edge/edge.py "$1" "$2"

