# 使用多阶段构建
FROM golang:1.21-alpine AS builder

# 设置工作目录
WORKDIR /app

# 复制go.mod和go.sum
COPY go.mod ./

# 下载依赖
RUN go mod download

# 复制源代码
COPY . .

# 构建应用
RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=7 go build -o subsmanager

# 使用轻量级基础镜像
FROM alpine:latest

# 安装必要的运行时依赖
RUN apk --no-cache add ca-certificates tzdata

# 设置工作目录
WORKDIR /app

# 从builder阶段复制二进制文件
COPY --from=builder /app/subsmanager .

# 创建数据目录
RUN mkdir -p /app/data

# 暴露端口
EXPOSE 3355

# 设置入口点
ENTRYPOINT ["./subsmanager"] 