#!/bin/bash

echo "🔍 Consul 服务发现测试脚本"
echo "=========================="

CONSUL_URL="http://127.0.0.1:8500"

echo ""
echo "1. 检查 Consul 服务状态..."
curl -s "$CONSUL_URL/v1/status/leader" | jq .

echo ""
echo "2. 查看所有注册的服务..."
curl -s "$CONSUL_URL/v1/catalog/services" | jq .

echo ""
echo "3. 查看 ai-service 详细信息..."
curl -s "$CONSUL_URL/v1/catalog/service/ai-service" | jq .

echo ""
echo "4. 查看 ai-service 健康检查状态..."
curl -s "$CONSUL_URL/v1/health/service/ai-service" | jq .

echo ""
echo "5. 测试服务发现 - 通过 Consul 获取服务地址..."
SERVICE_INFO=$(curl -s "$CONSUL_URL/v1/catalog/service/ai-service")
SERVICE_ADDRESS=$(echo $SERVICE_INFO | jq -r '.[0].ServiceAddress')
SERVICE_PORT=$(echo $SERVICE_INFO | jq -r '.[0].ServicePort')

echo "发现的服务地址: $SERVICE_ADDRESS:$SERVICE_PORT"

if [ "$SERVICE_ADDRESS" != "null" ] && [ "$SERVICE_PORT" != "null" ]; then
    echo ""
    echo "6. 通过发现的服务地址测试健康检查..."
    curl -s "http://$SERVICE_ADDRESS:8003/health" | jq .
else
    echo "❌ 无法获取服务地址信息"
fi

echo ""
echo "✅ Consul 服务发现测试完成！"
