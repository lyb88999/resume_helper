# 简历优化系统 - 设置指南

## 🚀 快速开始

### 前置要求
- Docker & Docker Compose
- Node.js 18+ & npm
- Go 1.21+ (可选，用于本地开发)

### 1. 启动后端服务

```bash
# 启动所有后端服务
./scripts/start-services.sh
```

这将启动以下服务：
- **用户服务** (user-service): http://localhost:8000
- **文件服务** (file-service): http://localhost:8001  
- **AI服务** (ai-service): http://localhost:8002
- **解析服务** (parser-service): http://localhost:8003
- **Consul**: http://localhost:8500
- **MySQL**: localhost:3307
- **Redis**: localhost:6379

### 2. 启动前端服务

```bash
# 启动前端开发服务器
./scripts/start-frontend.sh
```

前端将在 http://localhost:3000 启动

### 3. 验证服务状态

```bash
# 检查所有服务状态
docker-compose -f docker-compose.dev.yml ps

# 检查服务健康状态
curl http://localhost:8000/health  # 用户服务
curl http://localhost:8001/health  # 文件服务
curl http://localhost:8002/health  # AI服务
curl http://localhost:8003/health  # 解析服务
```

## 🏗️ 项目架构

### 后端服务
- **用户服务**: 用户注册、登录、信息管理
- **文件服务**: 文件上传、存储、管理
- **AI服务**: 简历分析、智能建议、知识检索
- **解析服务**: 简历解析、数据提取

### 前端应用
- **Vue 3 + TypeScript**: 现代化前端框架
- **Element Plus**: UI组件库
- **Pinia**: 状态管理
- **Vite**: 构建工具

## 🔧 开发配置

### 环境变量
前端通过 `vite.config.ts` 中的代理配置连接到后端服务：

```typescript
proxy: {
  '/v1/user': 'http://localhost:8000',    // 用户服务
  '/v1/resume': 'http://localhost:8000',  // 简历服务
  '/v1/file': 'http://localhost:8001',    // 文件服务
  '/v1/ai': 'http://localhost:8002',      // AI服务
  '/v1/parser': 'http://localhost:8003'   // 解析服务
}
```

### 端口配置
- 前端: 3000
- 用户服务: 8000 (HTTP), 9000 (gRPC)
- 文件服务: 8001 (HTTP), 9001 (gRPC)
- AI服务: 8002 (HTTP), 9002 (gRPC)
- 解析服务: 8003 (HTTP), 9003 (gRPC)

## 🐛 故障排除

### 常见问题

1. **服务启动失败**
   ```bash
   # 查看服务日志
   docker-compose -f docker-compose.dev.yml logs [service-name]
   
   # 重新构建并启动
   docker-compose -f docker-compose.dev.yml up --build
   ```

2. **端口冲突**
   - 检查端口是否被占用: `lsof -i :[port]`
   - 修改 `docker-compose.dev.yml` 中的端口映射

3. **数据库连接失败**
   - 确保MySQL容器正常运行
   - 检查数据库连接配置

4. **前端无法连接后端**
   - 检查代理配置是否正确
   - 确认后端服务已启动
   - 查看浏览器控制台错误信息

### 日志查看
```bash
# 查看所有服务日志
docker-compose -f docker-compose.dev.yml logs -f

# 查看特定服务日志
docker-compose -f docker-compose.dev.yml logs -f [service-name]
```

## 📚 API文档

### 用户服务 API
- `POST /v1/user/register` - 用户注册
- `POST /v1/user/login` - 用户登录
- `GET /v1/user/{id}` - 获取用户信息
- `PUT /v1/user/{id}` - 更新用户信息

### 简历服务 API
- `GET /v1/resume/list` - 获取简历列表
- `POST /v1/resume/upload` - 上传简历
- `POST /v1/resume/parse` - 解析简历
- `POST /v1/resume/analyze` - 分析简历

### AI服务 API
- `POST /v1/ai/chat` - AI智能问答
- `POST /v1/ai/resume/analyze` - AI简历分析
- `POST /v1/ai/knowledge/retrieve` - 知识检索
- `GET /v1/ai/health` - 健康检查

## 🚀 部署

### 生产环境
```bash
# 构建前端
cd frontend/web
npm run build

# 启动生产服务
docker-compose -f docker-compose.yml up -d
```

### 环境变量配置
生产环境需要配置以下环境变量：
- 数据库连接信息
- Redis连接信息
- 服务注册中心地址
- API密钥和认证信息

## 📞 支持

如果遇到问题，请：
1. 查看服务日志
2. 检查配置文件
3. 确认网络连接
4. 查看故障排除部分

---

**注意**: 这是一个开发环境配置，生产环境需要额外的安全配置和优化。

