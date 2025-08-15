FROM golang:1.24-alpine AS builder

WORKDIR /app

# 复制go mod文件
COPY go.mod go.sum ./
RUN go mod download

# 复制源代码
COPY . .

# 构建应用
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o ai-service ./cmd/ai-service

# 运行阶段
FROM alpine:latest

RUN apk --no-cache add ca-certificates curl

WORKDIR /root/

# 从构建阶段复制二进制文件
COPY --from=builder /app/ai-service .

# 复制配置文件
COPY --from=builder /app/configs/ ./configs/

# 暴露端口
EXPOSE 8002 9002

# 健康检查
HEALTHCHECK --interval=30s --timeout=10s --start-period=40s --retries=3 \
  CMD curl -f http://localhost:8002/health || exit 1

# 运行应用
CMD ["./ai-service"]
