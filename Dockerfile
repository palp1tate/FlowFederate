# 使用 python:3.10-slim 作为基础镜像
FROM python:3.10-slim

# 创建或覆盖 sources.list 文件，使用中科大的镜像源
RUN echo "deb http://mirrors.ustc.edu.cn/debian/ bullseye main contrib non-free" > /etc/apt/sources.list && \
    echo "deb http://mirrors.ustc.edu.cn/debian-security/ bullseye-security main contrib non-free" >> /etc/apt/sources.list

# 更新系统并安装必要的系统依赖
RUN apt-get update && apt-get install -y netcat-openbsd tzdata && rm -rf /var/lib/apt/lists/*

# 设置时区为上海时间（东八区）
ENV TZ=Asia/Shanghai
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

# 设置工作目录为/app
WORKDIR /app

# 将当前目录下的所有文件复制到工作目录下
COPY . .

# 创建并激活虚拟环境
RUN python -m venv venv
ENV PATH="/app/venv/bin:$PATH"

# 升级pip
RUN pip install --upgrade pip -i https://pypi.tuna.tsinghua.edu.cn/simple

# 安装 Python 依赖
RUN pip install --no-cache-dir -r requirements.txt -i https://pypi.tuna.tsinghua.edu.cn/simple

# 给 run.sh 设置为可执行权限
RUN chmod +x ./service/edge/run.sh ./service/client/run.sh
