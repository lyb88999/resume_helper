# ResumeOptim - 智能简历优化系统

## 📋 项目概述

ResumeOptim是一个基于AI的智能简历优化系统，采用微服务架构，为求职者提供专业、个性化的简历分析和改进建议。

## 🏗️ 项目结构

```
resumeOptim_claude/
├── backend/                    # 后端微服务
│   ├── services/              # 微服务目录
│   │   ├── user-service/      # 用户服务
│   │   ├── file-service/      # 文件服务
│   │   ├── parser-service/    # 解析服务
│   │   ├── ai-service/        # AI分析服务
│   │   ├── knowledge-service/ # 知识库服务
│   │   └── notification-service/ # 通知服务
│   ├── gateway/               # API网关
│   ├── shared/                # 共享代码
│   │   ├── pkg/              # 公共包
│   │   ├── proto/            # Protobuf定义
│   │   └── config/           # 配置模板
│   └── scripts/              # 构建和部署脚本
├── frontend/                  # 前端项目
│   ├── web/                  # Vue3 Web应用
│   └── mobile/               # 移动端应用（预留）
├── deployment/               # 部署配置
│   ├── docker/              # Docker配置
│   ├── k8s/                 # Kubernetes配置
│   └── docker-compose/      # 本地开发环境
├── scripts/                 # 项目脚本
├── configs/                 # 配置文件
├── docs/                    # 项目文档
└── tools/                   # 开发工具
```

## 🚀 快速开始

### 前置要求
- Go 1.21+
- Node.js 18+
- Docker & Docker Compose
- PostgreSQL 15+
- Redis 7+

### 本地开发环境搭建

```bash
# 1. 克隆项目
git clone https://github.com/your-org/resumeOptim_claude.git
cd resumeOptim_claude

# 2. 启动基础设施
make dev-infra

# 3. 启动后端服务
make dev-backend

# 4. 启动前端服务
make dev-frontend

# 5. 访问应用
open http://localhost:3000
```

## 📖 文档

详细文档请查看 [docs/README.md](./docs/README.md)

## 🛠️ 开发指南

### 代码规范
- Go代码遵循官方规范，使用gofmt和golint
- TypeScript代码使用ESLint和Prettier
- Git提交信息遵循Conventional Commits规范

### 分支管理
- `main`: 主分支，生产环境代码
- `develop`: 开发分支，集成最新功能
- `feature/*`: 功能分支
- `hotfix/*`: 紧急修复分支

## 📊 技术栈

### 后端
- **框架**: Go + Kratos v2
- **数据库**: PostgreSQL, MongoDB, Redis, Milvus
- **消息队列**: RabbitMQ
- **注册中心**: Consul

### 前端
- **框架**: Vue 3 + TypeScript
- **UI库**: Element Plus
- **状态管理**: Pinia
- **构建工具**: Vite

### AI/ML
- **编排框架**: Eino
- **大语言模型**: GPT-4, Claude-3.5
- **向量数据库**: Milvus

### 部署
- **容器化**: Docker + Kubernetes
- **CI/CD**: GitLab CI/CD
- **监控**: Prometheus + Grafana

## 📄 许可证

MIT License - 查看 [LICENSE](./LICENSE) 文件了解详情

---

**项目状态**: 🚧 开发中  
**当前版本**: v0.1.0  
**维护团队**: 全栈开发团队
