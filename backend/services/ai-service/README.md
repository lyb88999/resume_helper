# AI Service - 简历优化AI服务

基于Eino编排框架的智能简历分析和优化服务，提供简历解析、智能分析、优化建议和智能问答等功能。

## 功能特性

### 🔍 智能简历分析
- **多格式支持**: 支持PDF、Word、Markdown等格式简历解析
- **结构化提取**: 自动提取个人信息、教育背景、工作经历、项目经验等
- **多维度评分**: 完整性、清晰度、关键词匹配、格式规范、量化程度等
- **个性化建议**: 基于目标职位生成针对性优化建议

### 🤖 Eino编排能力
- **简历解析Chain**: 文档加载 → 内容解析 → 结构化提取
- **智能分析Graph**: 并行分析多个维度，提供综合评估
- **React Agent**: 智能问答助手，支持多轮对话
- **工具集成**: 职位匹配、模板推荐、市场调研等工具

### 💬 智能问答
- **上下文理解**: 基于简历内容的智能对话
- **知识库检索**: 融合行业知识和最佳实践
- **多轮对话**: 支持复杂问题的深入讨论
- **个性化回答**: 根据用户背景提供定制化建议

### 📊 数据存储与缓存
- **MySQL存储**: 分析结果和会话数据持久化
- **Redis缓存**: 提升响应速度和系统性能
- **会话管理**: 支持长期会话和上下文维护

## 技术架构

### 核心技术栈
- **框架**: Go + Kratos微服务框架
- **AI编排**: Eino框架（字节跳动开源）
- **数据库**: MySQL + Redis
- **服务注册**: Consul
- **协议**: gRPC + HTTP
- **依赖注入**: Wire

### Eino编排组件

#### 1. 简历解析Chain (ResumeParsingChain)
```go
文档加载 → 文档解析 → 结构化提取 → 数据验证
```

#### 2. 智能分析Graph (AnalysisGraph)
```go
知识检索 → 并行分析 → 评分计算 → 建议生成
    ↓
[完整性分析、清晰度分析、关键词分析、格式分析、量化分析]
```

#### 3. React Agent (ResumeOptimizeAgent)
```go
问题理解 → 工具选择 → 执行推理 → 生成回答
```

## 快速开始

### 环境要求
- Go 1.24+
- MySQL 8.0+
- Redis 6.0+
- Consul 1.15+

### 安装依赖
```bash
# 安装开发工具
make init

# 下载依赖
make deps
```

### 配置设置
```bash
# 复制并修改配置文件
cp configs/config.yaml configs/config.local.yaml

# 设置环境变量
export OPENAI_API_KEY=your-openai-api-key
```

### 生成代码
```bash
# 生成protobuf文件
make proto

# 生成Wire依赖注入代码
make wire

# 生成所有文件
make generate
```

### 运行服务
```bash
# 使用本地配置运行
make run-local

# 或者使用默认配置
make run
```

### Docker部署
```bash
# 构建镜像
make docker-build

# 运行容器
make docker-run
```

## API接口

### 1. 简历分析
```bash
POST /api/v1/ai/analyze
{
  "resume_id": "resume_123",
  "content": "简历内容...",
  "file_type": "pdf",
  "target_position": "软件工程师",
  "options": {
    "enable_completeness": true,
    "enable_clarity": true,
    "enable_keyword": true,
    "enable_format": true,
    "enable_quantification": true
  }
}
```

### 2. 生成建议
```bash
POST /api/v1/ai/suggestions
{
  "analysis_id": "analysis_123",
  "target_position": "软件工程师",
  "industry": "互联网",
  "options": {
    "max_suggestions": 5,
    "focus_area": "technical_skills",
    "experience_level": "senior"
  }
}
```

### 3. 智能问答
```bash
POST /api/v1/ai/chat
{
  "session_id": "session_123",
  "message": "如何优化我的简历？",
  "context": "resume_context",
  "options": {
    "use_resume_context": true,
    "use_knowledge_base": true,
    "language": "zh"
  }
}
```

### 4. 知识检索
```bash
POST /api/v1/ai/knowledge/retrieve
{
  "query": "软件工程师简历优化",
  "top_k": 10,
  "similarity_threshold": 0.7,
  "filters": ["resume_tips", "career_guide"]
}
```

### 5. 健康检查
```bash
GET /api/v1/ai/health
GET /health
GET /ai/status
```

## 配置说明

### AI模型配置
```yaml
ai:
  model:
    provider: openai          # 模型提供商
    api_key: ${OPENAI_API_KEY} # API密钥
    base_url: https://api.openai.com/v1
    model_name: gpt-4          # 模型名称
    max_tokens: 4096           # 最大token数
    temperature: 0.7           # 温度参数
    timeout_seconds: 60        # 超时时间
```

### Eino框架配置
```yaml
ai:
  eino:
    enable_tracing: true       # 启用链路追踪
    enable_caching: true       # 启用缓存
    max_concurrent: 10         # 最大并发数
    log_level: info           # 日志级别
```

### 向量数据库配置
```yaml
ai:
  vector:
    provider: milvus           # 向量数据库提供商
    address: 127.0.0.1:19530  # 数据库地址
    collection_name: resume_knowledge
    dimension: 1536           # 向量维度
    similarity_threshold: 0.7 # 相似度阈值
```

## 开发指南

### 项目结构
```
ai-service/
├── cmd/ai-service/          # 服务入口
├── configs/                 # 配置文件
├── api/ai/v1/              # API定义
├── internal/               # 内部模块
│   ├── biz/                # 业务逻辑层
│   ├── data/               # 数据访问层
│   ├── service/            # 服务层
│   ├── server/             # 服务器配置
│   ├── conf/               # 配置结构
│   └── eino/               # Eino编排组件
│       ├── schema.go       # 数据结构定义
│       ├── chains/         # Chain编排
│       ├── graphs/         # Graph编排
│       └── agents/         # Agent智能体
├── third_party/            # 第三方proto文件
├── Dockerfile             # Docker配置
├── Makefile              # 构建脚本
└── README.md             # 项目文档
```

### 开发流程

1. **添加新功能**
   ```bash
   # 1. 修改proto定义
   vim api/ai/v1/ai.proto
   
   # 2. 生成代码
   make proto
   
   # 3. 实现业务逻辑
   vim internal/biz/ai.go
   
   # 4. 更新服务层
   vim internal/service/ai.go
   
   # 5. 重新生成依赖注入
   make wire
   ```

2. **添加Eino组件**
   ```bash
   # 1. 定义新的Chain/Graph/Agent
   vim internal/eino/chains/new_chain.go
   
   # 2. 集成到业务逻辑
   vim internal/biz/ai.go
   
   # 3. 测试组件
   make test
   ```

3. **调试和测试**
   ```bash
   # 运行测试
   make test
   
   # 代码检查
   make lint
   
   # 本地调试
   make run-local
   ```

### 性能优化

#### 1. 缓存策略
- **分析结果缓存**: Redis缓存1小时
- **会话上下文缓存**: Redis缓存30分钟
- **知识库缓存**: 内存+Redis两级缓存

#### 2. 并发优化
- **并行分析**: Graph组件支持并行分析
- **连接池**: 数据库和Redis连接池
- **限流控制**: 防止AI API调用过载

#### 3. 监控告警
- **链路追踪**: 支持Jaeger/Zipkin
- **指标监控**: Prometheus指标
- **日志聚合**: 结构化日志输出

## 部署运维

### Docker Compose部署
```yaml
version: '3.8'
services:
  ai-service:
    build: .
    ports:
      - "8003:8003"
      - "9003:9003"
    environment:
      - OPENAI_API_KEY=${OPENAI_API_KEY}
    depends_on:
      - mysql
      - redis
      - consul
```

### Kubernetes部署
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: ai-service
spec:
  replicas: 3
  selector:
    matchLabels:
      app: ai-service
  template:
    spec:
      containers:
      - name: ai-service
        image: ai-service:latest
        ports:
        - containerPort: 8003
        - containerPort: 9003
```

### 健康检查
```bash
# 基础健康检查
curl http://localhost:8003/health

# AI服务状态检查
curl http://localhost:8003/ai/status

# 详细组件状态
curl http://localhost:8003/api/v1/ai/health
```

## 常见问题

### Q: 如何配置自定义AI模型？
A: 修改配置文件中的`ai.model`部分，支持OpenAI、Claude等模型。

### Q: 如何扩展新的分析维度？
A: 在`AnalysisGraph`中添加新的分析节点，并更新评分计算逻辑。

### Q: 如何优化向量检索性能？
A: 调整`similarity_threshold`参数，使用更高效的向量数据库。

### Q: 如何处理大文件解析？
A: 配置`max_file_size`限制，使用分块处理策略。

## 贡献指南

1. Fork 项目
2. 创建功能分支 (`git checkout -b feature/AmazingFeature`)
3. 提交修改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 开启 Pull Request

## 许可证

本项目采用 MIT 许可证 - 查看 [LICENSE](LICENSE) 文件了解详情。

## 更新日志

### v1.0.0 (2024-01-01)
- ✨ 初始版本发布
- 🔥 支持基于Eino的简历分析
- 🤖 集成React Agent智能问答
- 📊 完整的评分和建议系统
- 🚀 Docker和K8s部署支持

