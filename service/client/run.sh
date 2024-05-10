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

# 设置 PYTHONPATH 环境变量
export PYTHONPATH=./

# 执行 python 程序
python ./service/client/client.py "$1" "$2"

