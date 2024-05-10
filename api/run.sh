#!/bin/sh

# 环境变量 MYSQL_HOST 和 MYSQL_PORT 应该设置为 MySQL 服务器的主机名和端口
host="${MYSQL_HOST:-mysql}"
port="${MYSQL_PORT:-3306}"

# 尝试连接到 MySQL
while ! nc -z "$host" "$port"; do
  echo "Waiting for MySQL at $host:$port..."
  sleep 1
done

echo "MySQL is up..."

# 执行 main 程序
./main -p 9090

