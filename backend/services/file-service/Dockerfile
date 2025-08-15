FROM golang:1.24-alpine AS builder

# 设置工作目录
WORKDIR /app

# 安装必要的包
RUN apk add --no-cache git

# 复制go模块文件
COPY go.mod ./
COPY go.sum ./

# 下载依赖
RUN go mod download

# 复制源代码
COPY . .

# 构建应用
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o file-service ./cmd/file-service

FROM alpine:latest

# 安装ca-certificates和curl用于健康检查
RUN apk --no-cache add ca-certificates curl

WORKDIR /app

# 从builder阶段复制二进制文件
COPY --from=builder /app/file-service .

# 复制配置文件
COPY configs/ ./configs/

# 创建上传目录
RUN mkdir -p /app/uploads

# 暴露端口
EXPOSE 8001 9001

# 健康检查
HEALTHCHECK --interval=30s --timeout=10s --start-period=40s --retries=3 \
  CMD curl -f http://localhost:8001/health || exit 1

# 启动应用
CMD ["./file-service", "-conf", "./configs/config.yaml"]