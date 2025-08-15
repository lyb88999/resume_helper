#!/bin/bash

echo "🧪 前端API调用测试脚本"
echo "======================"

# 检查前端服务是否运行
echo ""
echo "1. 检查前端服务状态..."
if curl -s http://localhost:3000 > /dev/null; then
    echo "✅ 前端服务运行在 http://localhost:3000"
else
    echo "❌ 前端服务未运行，请先启动前端服务"
    echo "   在 frontend/web 目录下运行: npm run dev"
    exit 1
fi

# 测试用户服务
echo ""
echo "2. 测试用户服务 (通过代理 /v1/user -> localhost:8001)..."
if curl -s http://localhost:8001/health > /dev/null; then
    echo "✅ 用户服务健康检查通过"
else
    echo "❌ 用户服务未运行或无法访问"
fi

# 测试解析服务
echo ""
echo "3. 测试解析服务 (通过代理 /v1/resume -> localhost:8002)..."
if curl -s http://localhost:8002/health > /dev/null; then
    echo "✅ 解析服务健康检查通过"
else
    echo "❌ 解析服务未运行或无法访问"
fi

# 测试AI服务
echo ""
echo "4. 测试AI服务 (通过代理 /v1/ai -> localhost:8003)..."
if curl -s http://localhost:8003/health > /dev/null; then
    echo "✅ AI服务健康检查通过"
else
    echo "❌ AI服务未运行或无法访问"
fi

# 测试前端代理
echo ""
echo "5. 测试前端代理配置..."
echo "   测试 /v1/ai/status 代理到 ai-service..."

PROXY_TEST=$(curl -s -o /dev/null -w "%{http_code}" http://localhost:3000/v1/ai/status)
if [ "$PROXY_TEST" = "200" ]; then
    echo "✅ 前端代理配置正确，/v1/ai 成功代理到 ai-service"
else
    echo "❌ 前端代理配置有问题，状态码: $PROXY_TEST"
fi

echo ""
echo "6. 测试前端页面访问..."
FRONTEND_TEST=$(curl -s -o /dev/null -w "%{http_code}" http://localhost:3000)
if [ "$FRONTEND_TEST" = "200" ]; then
    echo "✅ 前端页面可以正常访问"
else
    echo "❌ 前端页面访问失败，状态码: $FRONTEND_TEST"
fi

echo ""
echo "📋 服务状态总结:"
echo "=================="
echo "前端服务: http://localhost:3000"
echo "用户服务: http://localhost:8001"
echo "解析服务: http://localhost:8002"
echo "AI服务:   http://localhost:8003"

echo ""
echo "🔧 如果服务未运行，请启动相应服务:"
echo "   用户服务: cd backend/services/user-service && ./bin/user-service -conf configs/config.local.yaml"
echo "   解析服务: cd backend/services/parser-service && ./bin/parser-service -conf configs/config.local.yaml"
echo "   AI服务:   cd backend/services/ai-service && ./bin/ai-service -conf configs/config.local.yaml"

echo ""
echo "✅ 前端API测试完成！"

