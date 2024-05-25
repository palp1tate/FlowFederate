#!/bin/sh

# 尝试连接到Nacos

while ! nc -z nacos 8848; do
  echo "Waiting for Nacos..."
  sleep 1
done

echo "Nacos is up..."

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

# 执行 main 程序
./main "$1" "$2"
