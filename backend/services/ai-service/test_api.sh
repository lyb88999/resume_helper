#!/bin/bash

echo "🚀 AI Service API 测试脚本"
echo "=========================="

BASE_URL="http://localhost:8003"

echo ""
echo "1. 测试健康检查接口..."
curl -s -X GET "$BASE_URL/health" | jq .

echo ""
echo "2. 测试AI服务状态..."
curl -s -X GET "$BASE_URL/ai/status" | jq .

echo ""
echo "3. 测试简历分析接口..."
curl -s -X POST "$BASE_URL/api/v1/ai/analyze" \
  -H "Content-Type: application/json" \
  -d '{
    "resume_id": "test_001",
    "target_position": "软件工程师",
    "file_path": "",
    "options": {
      "enable_completeness": true,
      "enable_clarity": true,
      "enable_keyword": true,
      "enable_format": true,
      "enable_quantification": true
    }
  }' | jq .

echo ""
echo "4. 测试AI聊天接口..."
curl -s -X POST "$BASE_URL/api/v1/ai/chat" \
  -H "Content-Type: application/json" \
  -d '{
    "message": "如何优化我的简历？",
    "session_id": "test_session_001",
    "context": "我是一名应届毕业生，想找软件工程师的工作"
  }' | jq .

echo ""
echo "✅ 所有API测试完成！"
