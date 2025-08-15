#!/bin/bash

echo "🚀 启动简历优化系统服务..."

# 检查Docker是否运行
if ! docker info > /dev/null 2>&1; then
    echo "❌ Docker未运行，请先启动Docker"
    exit 1
fi

# 停止并删除现有容器
echo "🔄 清理现有容器..."
docker-compose -f docker-compose.dev.yml down

# 构建并启动服务
echo "🔨 构建并启动服务..."
docker-compose -f docker-compose.dev.yml up --build -d

# 等待服务启动
echo "⏳ 等待服务启动..."
sleep 30

# 检查服务状态
echo "📊 检查服务状态..."
docker-compose -f docker-compose.dev.yml ps

# 检查服务健康状态
echo "🏥 检查服务健康状态..."
services=("user-service:8000" "file-service:8001" "ai-service:8002" "parser-service:8003")

for service in "${services[@]}"; do
    service_name=$(echo $service | cut -d: -f1)
    port=$(echo $service | cut -d: -f2)
    
    echo "检查 $service_name..."
    if curl -f "http://localhost:$port/health" > /dev/null 2>&1; then
        echo "✅ $service_name 运行正常"
    else
        echo "❌ $service_name 启动失败"
    fi
done

echo ""
echo "🎉 服务启动完成！"
echo "📱 前端地址: http://localhost:3000"
echo "🔧 用户服务: http://localhost:8000"
echo "📁 文件服务: http://localhost:8001"
echo "🤖 AI服务: http://localhost:8002"
echo "📄 解析服务: http://localhost:8003"
echo "🗄️  Consul: http://localhost:8500"
echo "🐘 MySQL: localhost:3307"
echo "🔴 Redis: localhost:6379"
