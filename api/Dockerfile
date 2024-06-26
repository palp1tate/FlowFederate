# 使用 golang:1.22 作为构建阶段的基础镜像
FROM golang:1.22 AS builder

# 设置工作目录为/app
WORKDIR /app

# 将当前目录下的所有文件复制到工作目录下
COPY . .

# 设置 Go 的代理服务器
RUN go env -w GOPROXY=https://goproxy.cn,direct

# 下载项目依赖
RUN go mod download

# 构建项目生成可执行文件
RUN go build -o main ./api/main.go

# 使用 ubuntu:22.04 作为最终阶段的基础镜像
FROM ubuntu:22.04

# 替换默认Ubuntu源为中国科学技术大学的源
# 创建或覆盖 sources.list 文件，使用中科大的 Ubuntu 镜像源
RUN echo "deb http://mirrors.ustc.edu.cn/ubuntu/ focal main restricted universe multiverse" > /etc/apt/sources.list && \
    echo "deb http://mirrors.ustc.edu.cn/ubuntu/ focal-security main restricted universe multiverse" >> /etc/apt/sources.list && \
    echo "deb http://mirrors.ustc.edu.cn/ubuntu/ focal-updates main restricted universe multiverse" >> /etc/apt/sources.list && \
    echo "deb http://mirrors.ustc.edu.cn/ubuntu/ focal-proposed main restricted universe multiverse" >> /etc/apt/sources.list && \
    echo "deb http://mirrors.ustc.edu.cn/ubuntu/ focal-backports main restricted universe multiverse" >> /etc/apt/sources.list

# 更新系统并安装必要的系统依赖
RUN apt-get update && apt-get install -y netcat-openbsd tzdata && rm -rf /var/lib/apt/lists/*

# 设置时区为上海时间（东八区）
ENV TZ=Asia/Shanghai
RUN ln -snf /usr/share/zoneinfo/"$TZ" /etc/localtime && echo "$TZ" > /etc/timezone

# 设置工作目录为/app
WORKDIR /app

# 将构建阶段的可执行文件和脚本复制到工作目录下
COPY --from=builder /app/main /app/api/run.sh ./

# 将构建阶段的内部文件复制到工作目录的内部文件夹中
COPY --from=builder /app/internal ./internal/

# 给 run.sh 和 main 设置为可执行权限
RUN chmod +x ./run.sh ./main