# 构建阶段
FROM golang:1.24-alpine AS builder

# 设置工作目录
WORKDIR /app

# 安装必要的包
RUN apk add --no-cache git ca-certificates

# 复制 go mod 文件
COPY go.mod ./
COPY go.sum ./

# 复制源代码
COPY . .

# 构建应用
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main cmd/main.go cmd/wire_gen.go

# 运行阶段
FROM alpine:latest

# 安装 ca-certificates 和 curl (用于健康检查)
RUN apk --no-cache add ca-certificates curl

WORKDIR /root/

# 从构建阶段复制二进制文件
COPY --from=builder /app/main .

# 复制配置文件
COPY configs/ ./configs/

# 暴露端口
EXPOSE 8000 9000

# 健康检查
HEALTHCHECK --interval=30s --timeout=10s --start-period=40s --retries=3 \
  CMD curl -f http://localhost:8000/health || exit 1

# 运行应用
CMD ["./main", "-conf", "./configs/user-service.yaml"]
