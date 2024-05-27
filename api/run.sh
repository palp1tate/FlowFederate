#!/bin/sh

# 尝试连接到 MySQL
while ! nc -z mysql 3306; do
  echo "Waiting for MySQL..."
  sleep 1
done

echo "MySQL is up..."

# 尝试连接到 Redis
while ! nc -z redis 6379; do
  echo "Waiting for Redis..."
  sleep 1
done

echo "Redis is up..."

# 尝试连接到 Consul
while ! nc -z consul 8500; do
  echo "Waiting for Consul..."
  sleep 1
done

echo "Consul is up..."

# 执行 main 程序
./main "$1" "$2"
