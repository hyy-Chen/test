# 使用基础镜像
FROM alpine:latest

# 设置工作目录
WORKDIR /app

# 复制可执行文件到容器中
COPY main /app/main

# 暴露端口（如果需要）
EXPOSE 8080

# 启动命令
CMD ["./main"]
