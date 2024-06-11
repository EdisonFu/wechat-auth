FROM ubuntu:latest
LABEL authors="fuzd"

# 设置基础镜像
FROM golang:1.21-alpine

# 设置工作目录
WORKDIR /app

RUN go env -w GOPROXY="https://goproxy.cn,direct"

COPY go.mod go.sum ./

RUN go mod download

# 复制项目文件到容器中
COPY . .

# 构建Go应用
RUN go build -o azgw

# 设置容器启动命令
CMD ["./azgw"]