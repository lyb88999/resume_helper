# 简历优化系统 - 开发指南

## 🚀 快速开始

这是一个基于 **Kratos + Eino + Vue.js** 构建的智能简历优化系统。

### 📋 前置要求

- **Go** 1.21+
- **Node.js** 18+
- **Docker** & **Docker Compose**
- **Git**

### ⚡ 一键启动

```bash
# 1. 克隆项目
git clone https://github.com/liyubo06/resumeOptim_claude.git
cd resumeOptim_claude

# 2. 初始化开发环境
make setup

# 3. 配置环境变量
cp .env.example .env
# 编辑 .env 文件，设置你的 OPENAI_API_KEY

# 4. 启动系统
make quick-start
```

系统启动后，访问：
- 🎨 **前端应用**: http://localhost:3000
- 🔌 **API网关**: http://localhost:8080
- 📊 **监控面板**: http://localhost:3001 (admin/admin123)

## 🏗️ 项目架构

### 技术栈

| 组件 | 技术选型 | 说明 |
|------|----------|------|
| 🎨 **前端** | Vue.js 3 + TypeScript + Element Plus | 现代化响应式界面 |
| 🔧 **后端** | Kratos (Go) + Eino (AI编排) | 微服务架构 + AI工作流 |
| 🗄️ **数据库** | MySQL + Redis + Milvus | 关系型 + 缓存 + 向量数据库 |
| 🐳 **部署** | Docker + Kubernetes | 容器化云原生部署 |
| 📊 **监控** | Prometheus + Grafana + Jaeger | 可观测性体系 |

### 系统架构

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Vue.js 前端   │───▶│   API Gateway   │───▶│   微服务集群     │
│   - Element UI  │    │   - 路由转发     │    │   - 用户服务     │
│   - Pinia状态   │    │   - 认证鉴权     │    │   - 解析服务     │
│   - TypeScript  │    │   - 限流熔断     │    │   - AI分析服务   │
└─────────────────┘    └─────────────────┘    │   - 知识库服务   │
                                             │   - 文件服务     │
┌─────────────────┐    ┌─────────────────┐    └─────────────────┘
│   Eino AI引擎   │───▶│   基础设施层     │
│   - Chain编排   │    │   - MySQL数据库  │         ┌──────────┐
│   - Graph分析   │    │   - Redis缓存    │────────▶│ 监控告警  │
│   - Agent智能体 │    │   - Milvus向量库 │         │ Grafana  │
└─────────────────┘    │   - Kafka消息队列│         └──────────┘
                       └─────────────────┘
```

## 🔧 开发环境

### 目录结构

```
resumeOptim_claude/
├── backend/                    # 后端服务
│   ├── gateway/               # API网关
│   ├── services/              # 微服务集群
│   │   ├── user-service/      # 用户服务
│   │   ├── parser-service/    # 解析服务
│   │   ├── ai-service/        # AI分析服务
│   │   ├── knowledge-service/ # 知识库服务
│   │   └── file-service/      # 文件服务
│   └── shared/                # 共享代码
│       ├── config/            # 配置管理
│       ├── pkg/               # 工具包
│       └── proto/             # gRPC定义
├── frontend/                   # 前端应用
│   └── web/                   # Vue.js Web应用
├── deployment/                 # 部署配置
│   ├── docker-compose/        # Docker编排
│   └── k8s/                   # Kubernetes配置
├── doc/                       # 项目文档
├── configs/                   # 配置文件
└── scripts/                   # 开发脚本
```

### 🚀 启动方式

#### 方式一：完整Docker启动

```bash
# 启动所有服务（推荐新手）
make quick-start

# 查看服务状态
make docker-logs

# 停止所有服务
make stop
```

#### 方式二：混合开发模式

```bash
# 启动基础设施（数据库、缓存等）
make setup-infra

# 启动后端服务（本地运行）
make dev-backend

# 启动前端服务（本地运行）
make dev-frontend
```

#### 方式三：单独启动

```bash
# 仅启动基础设施
make setup-infra

# 手动启动服务
cd backend/services/user-service
go run cmd/main.go -conf ../../../configs/user-service.yaml

# 前端开发
cd frontend/web
npm run dev
```

## 🤖 Eino AI 工作流

### 核心概念

本系统基于字节跳动开源的 **Eino AI编排框架** 构建智能分析能力：

```go
// Chain编排 - 简历解析链
uploadFile → DocumentLoader → DocumentParser → 
ChatModel → StructureExtraction → ValidationResult

// Graph编排 - 并行分析图
KnowledgeRetrieval → ParallelAnalysis → ScoreCalculation → SuggestionGeneration
                     ├─ CompletenessAnalysis
                     ├─ ClarityAnalysis  
                     ├─ KeywordMatching
                     ├─ FormatValidation
                     └─ QuantificationAnalysis

// Agent应用 - 智能问答
UserInput → ReactAgent → ToolCalling → ResultIntegration → SmartReply
```

### AI工作流配置

```go
// 简历分析工作流示例
resumeAnalysisWorkflow := &compose.Graph{
    Nodes: []Node{
        {Name: "knowledge_retrieval", Type: "retriever"},
        {Name: "parallel_analysis", Type: "parallel"},
        {Name: "score_calculation", Type: "calculator"},
        {Name: "suggestion_generation", Type: "llm"},
    },
    Edges: []Edge{
        {From: "knowledge_retrieval", To: "parallel_analysis"},
        {From: "parallel_analysis", To: "score_calculation"},
        {From: "score_calculation", To: "suggestion_generation"},
    },
}
```

## 📊 开发工具

### Make命令

```bash
# 查看所有命令
make help

# 环境管理
make setup          # 初始化开发环境
make deps           # 安装依赖
make clean          # 清理构建产物

# 开发调试
make dev            # 启动开发环境
make test           # 运行测试
make lint           # 代码检查
make fmt            # 代码格式化

# 构建部署
make build          # 构建所有服务
make docker-build   # 构建Docker镜像
make docker-up      # 启动Docker服务

# 代码生成
make proto          # 生成protobuf代码
make wire           # 生成依赖注入代码
```

### 开发调试

```bash
# 查看服务日志
docker-compose -f deployment/docker-compose/docker-compose.yml logs -f [service_name]

# 进入容器调试
docker exec -it resumeoptim_claude_mysql_1 mysql -u root -p

# 查看系统状态
curl http://localhost:8080/health
```

## 🔌 API开发

### gRPC服务

```proto
// 用户服务示例
service UserService {
  rpc Register(RegisterRequest) returns (RegisterResponse);
  rpc Login(LoginRequest) returns (LoginResponse);
  rpc GetUser(GetUserRequest) returns (GetUserResponse);
}
```

### RESTful API

```bash
# 用户注册
POST /api/v1/users/register
{
  "username": "demo",
  "email": "demo@example.com", 
  "password": "password123"
}

# 简历上传
POST /api/v1/resumes/upload
Content-Type: multipart/form-data
file: resume.pdf

# 简历分析
POST /api/v1/analysis/analyze
{
  "resume_id": 123,
  "target_position": "软件工程师",
  "industry": "互联网"
}
```

## 🎨 前端开发

### 技术栈

- **Vue 3**: Composition API + `<script setup>`
- **TypeScript**: 类型安全
- **Element Plus**: UI组件库
- **Pinia**: 状态管理
- **Vite**: 构建工具

### 开发规范

```typescript
// 组件示例
<template>
  <div class="resume-upload">
    <el-upload
      :action="uploadUrl"
      :on-success="handleSuccess"
      :before-upload="beforeUpload"
    >
      <el-button type="primary">上传简历</el-button>
    </el-upload>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useResumeStore } from '@/stores/resume'
import type { UploadFile } from 'element-plus'

const resumeStore = useResumeStore()
const uploadUrl = ref('/api/v1/resumes/upload')

const beforeUpload = (file: UploadFile) => {
  // 文件类型和大小验证
  return true
}

const handleSuccess = (response: any) => {
  // 处理上传成功
  resumeStore.addResume(response.data)
}
</script>
```

### 状态管理

```typescript
// Pinia Store示例
export const useResumeStore = defineStore('resume', () => {
  const resumes = ref<Resume[]>([])
  const currentResume = ref<Resume | null>(null)
  
  const uploadResume = async (file: File) => {
    const response = await resumeApi.upload(file)
    resumes.value.push(response.data)
  }
  
  return { resumes, currentResume, uploadResume }
})
```

## 🧪 测试

### 后端测试

```bash
# 运行所有测试
make test-backend

# 运行特定服务测试
cd backend/services/user-service
go test -v ./...

# 生成测试覆盖率
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### 前端测试

```bash
# 运行前端测试
make test-frontend

# 单元测试
cd frontend/web
npm run test:unit

# E2E测试
npm run test:e2e
```

## 📊 监控与日志

### 服务监控

- **Grafana**: http://localhost:3001 (admin/admin123)
- **Prometheus**: http://localhost:9090
- **Jaeger**: http://localhost:16686
- **Consul**: http://localhost:8500

### 日志查看

```bash
# 查看所有服务日志
make docker-logs

# 查看特定服务日志
docker-compose -f deployment/docker-compose/docker-compose.yml logs -f user-service

# 实时监控
tail -f logs/application.log
```

## 🔧 故障排除

### 常见问题

#### 1. 服务启动失败

```bash
# 检查端口占用
lsof -i :8080

# 检查Docker服务状态
docker-compose ps

# 重启服务
make stop && make quick-start
```

#### 2. 数据库连接失败

```bash
# 检查MySQL状态
docker exec -it mysql mysql -u root -p -e "SHOW DATABASES;"

# 重建数据库
docker-compose down -v
docker-compose up -d mysql
```

#### 3. 前端构建失败

```bash
# 清理缓存
cd frontend/web
rm -rf node_modules package-lock.json
npm install

# 检查Node版本
node --version  # 需要 18+
```

#### 4. AI服务异常

```bash
# 检查OpenAI API Key
echo $OPENAI_API_KEY

# 查看AI服务日志
docker-compose logs ai-service

# 重启AI服务
docker-compose restart ai-service
```

### 性能优化

```bash
# 系统性能监控
make monitor

# 压力测试
make benchmark

# 内存使用分析
go tool pprof http://localhost:8080/debug/pprof/heap
```

## 🤝 贡献指南

### 开发流程

1. **Fork项目** → 创建特性分支
2. **本地开发** → 编写代码和测试
3. **提交PR** → 代码审查和合并

### 代码规范

```bash
# 代码格式化
make fmt

# 代码检查
make lint

# 安全检查
make security

# 质量检查
make quality
```

### 提交规范

```bash
# 提交格式
git commit -m "feat(user): 添加用户注册功能"
git commit -m "fix(parser): 修复PDF解析错误"
git commit -m "docs(api): 更新API文档"
```

## 📚 相关资源

### 官方文档

- [Eino框架文档](https://www.cloudwego.io/zh/docs/eino/quick_start/)
- [Kratos框架文档](https://go-kratos.dev/)
- [Vue.js文档](https://vuejs.org/)
- [Element Plus文档](https://element-plus.org/)

### 学习资源

- [Go语言学习](https://golang.org/doc/)
- [TypeScript教程](https://www.typescriptlang.org/docs/)
- [Docker教程](https://docs.docker.com/)
- [Kubernetes教程](https://kubernetes.io/docs/)

---

## 📞 技术支持

如遇问题，请：

1. 查看 [FAQ文档](./doc/FAQ.md)
2. 搜索 [Issues](https://github.com/liyubo06/resumeOptim_claude/issues)
3. 提交新Issue或联系维护者

**Happy Coding! 🚀**
