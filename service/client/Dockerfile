# 使用 python:3.10-slim 作为基础镜像
FROM python:3.10-slim

# 设置镜像源为中科大
RUN sed -i 's/deb.debian.org/mirrors.ustc.edu.cn/g' /etc/apt/sources.list && \
    sed -i 's/security.debian.org/mirrors.ustc.edu.cn/g' /etc/apt/sources.list

# 更新系统并安装必要的系统依赖
RUN apt-get update && apt-get install -y netcat && rm -rf /var/lib/apt/lists/*

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
RUN chmod +x ./service/client/run.sh

# 暴露端口
EXPOSE 5001

# 当容器启动时运行 run.sh
CMD ["./service/client/run.sh","-p","5001"]
