#!/bin/bash

# 简历优化系统开发环境设置脚本
set -e

echo "🚀 开始设置简历优化系统开发环境..."

# 检查必要工具
check_prerequisites() {
    echo "📋 检查前置条件..."
    
    # 检查Go
    if ! command -v go &> /dev/null; then
        echo "❌ Go未安装，请先安装Go 1.21+"
        exit 1
    fi
    
    # 检查Node.js
    if ! command -v node &> /dev/null; then
        echo "❌ Node.js未安装，请先安装Node.js 18+"
        exit 1
    fi
    
    # 检查Docker
    if ! command -v docker &> /dev/null; then
        echo "❌ Docker未安装，请先安装Docker"
        exit 1
    fi
    
    # 检查Docker Compose
    if ! command -v docker-compose &> /dev/null; then
        echo "❌ Docker Compose未安装，请先安装Docker Compose"
        exit 1
    fi
    
    echo "✅ 前置条件检查通过"
}

# 初始化后端项目
setup_backend() {
    echo "🔧 设置后端项目..."
    
    # 初始化Go模块
    echo "📦 初始化Go模块..."
    go mod tidy
    
    # 安装wire代码生成工具
    echo "🛠️  安装代码生成工具..."
    go install github.com/google/wire/cmd/wire@latest
    go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
    
    # 生成proto文件
    echo "📄 生成protobuf文件..."
    find backend/shared/proto -name "*.proto" -exec protoc \
        --go_out=. \
        --go_opt=paths=source_relative \
        --go-grpc_out=. \
        --go-grpc_opt=paths=source_relative \
        {} \;
    
    echo "✅ 后端项目设置完成"
}

# 初始化前端项目
setup_frontend() {
    echo "🎨 设置前端项目..."
    
    cd frontend/web
    
    # 安装依赖
    echo "📦 安装前端依赖..."
    npm install
    
    # 生成类型定义
    echo "🔧 生成类型定义..."
    npm run type-check
    
    cd ../..
    
    echo "✅ 前端项目设置完成"
}

# 启动基础设施服务
start_infrastructure() {
    echo "🏗️  启动基础设施服务..."
    
    cd deployment/docker-compose
    
    # 启动基础服务
    echo "🚀 启动数据库和缓存服务..."
    docker-compose up -d mysql redis consul milvus etcd minio
    
    # 等待服务启动
    echo "⏳ 等待服务启动完成..."
    sleep 30
    
    # 检查服务状态
    echo "🔍 检查服务状态..."
    docker-compose ps
    
    cd ../..
    
    echo "✅ 基础设施服务启动完成"
}

# 创建配置文件
create_configs() {
    echo "📝 创建配置文件..."
    
    # 创建配置目录
    mkdir -p configs
    
    # 创建开发环境配置
    cat > configs/config.yaml << EOF
server:
  http:
    port: 8080
    timeout: 30
  grpc:
    port: 9000
    timeout: 30

database:
  driver: mysql
  host: localhost
  port: 3306
  username: root
  password: resume123
  database: resume_optim
  max_open_conns: 100
  max_idle_conns: 10

redis:
  host: localhost
  port: 6379
  password: redis123
  db: 0

milvus:
  host: localhost
  port: 19530
  username: ""
  password: ""
  database: default
  collection: resume_knowledge

eino:
  model_provider: openai
  openai:
    api_key: \${OPENAI_API_KEY}
    base_url: https://api.openai.com/v1
    model: gpt-4
    temperature: 0.1
    max_tokens: 2000
  embeddings:
    provider: openai
    model: text-embedding-ada-002
    api_key: \${OPENAI_API_KEY}
  workflows:
    resume_parsing_timeout: 120
    analysis_timeout: 180
    knowledge_retrieval: 20
    max_concurrency: 5

log:
  level: debug
  encoding: console
  output_path: stdout

tracing:
  enabled: true
  service_name: resume-optim
  endpoint: http://localhost:14268/api/traces
EOF
    
    echo "✅ 配置文件创建完成"
}

# 创建环境变量文件
create_env_file() {
    echo "🔑 创建环境变量文件..."
    
    cat > .env << 'EOF'
# OpenAI API密钥 (必须设置)
OPENAI_API_KEY=your_openai_api_key_here

# 数据库配置
MYSQL_ROOT_PASSWORD=resume123
MYSQL_DATABASE=resume_optim
MYSQL_USER=resume_user
MYSQL_PASSWORD=resume123

# Redis配置
REDIS_PASSWORD=redis123

# MinIO配置
MINIO_ACCESS_KEY=minioadmin
MINIO_SECRET_KEY=minioadmin

# 应用环境
NODE_ENV=development
GO_ENV=development

# 服务端口配置
FRONTEND_PORT=3000
GATEWAY_PORT=8080
USER_SERVICE_PORT=8001
PARSER_SERVICE_PORT=8002
AI_SERVICE_PORT=8003
KNOWLEDGE_SERVICE_PORT=8004
FILE_SERVICE_PORT=8005
EOF
    
    echo "⚠️  请编辑 .env 文件，设置你的 OpenAI API Key"
    echo "✅ 环境变量文件创建完成"
}

# 初始化数据库
init_database() {
    echo "🗄️  初始化数据库..."
    
    # 等待MySQL启动
    echo "⏳ 等待MySQL启动..."
    sleep 20
    
    # 创建数据库表
    echo "📋 创建数据库表..."
    cat > configs/mysql/init.sql << 'EOF'
-- 创建数据库
CREATE DATABASE IF NOT EXISTS resume_optim CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
USE resume_optim;

-- 用户表
CREATE TABLE IF NOT EXISTS users (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    full_name VARCHAR(100),
    avatar_url VARCHAR(255),
    status TINYINT DEFAULT 1,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    INDEX idx_username (username),
    INDEX idx_email (email),
    INDEX idx_status (status)
);

-- 简历表
CREATE TABLE IF NOT EXISTS resumes (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    user_id BIGINT NOT NULL,
    title VARCHAR(200) NOT NULL,
    file_path VARCHAR(500),
    file_type ENUM('pdf', 'markdown') NOT NULL,
    file_size BIGINT,
    parsed_content JSON,
    status TINYINT DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    INDEX idx_user_id (user_id),
    INDEX idx_status (status),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- 分析结果表
CREATE TABLE IF NOT EXISTS analysis_results (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    resume_id BIGINT NOT NULL,
    overall_score DECIMAL(3,1),
    completeness_score DECIMAL(3,1),
    clarity_score DECIMAL(3,1),
    keyword_score DECIMAL(3,1),
    format_score DECIMAL(3,1),
    quantification_score DECIMAL(3,1),
    suggestions JSON,
    analysis_version VARCHAR(20),
    target_position VARCHAR(100),
    industry VARCHAR(50),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    INDEX idx_resume_id (resume_id),
    FOREIGN KEY (resume_id) REFERENCES resumes(id) ON DELETE CASCADE
);

-- 知识库表
CREATE TABLE IF NOT EXISTS knowledge_base (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    title VARCHAR(200) NOT NULL,
    content TEXT NOT NULL,
    category VARCHAR(50),
    tags JSON,
    file_path VARCHAR(500),
    vector_id VARCHAR(100),
    status TINYINT DEFAULT 1,
    created_by BIGINT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    INDEX idx_category (category),
    INDEX idx_status (status),
    INDEX idx_created_by (created_by),
    FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE SET NULL
);

-- 知识块表
CREATE TABLE IF NOT EXISTS knowledge_chunks (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    knowledge_id BIGINT NOT NULL,
    chunk_text TEXT NOT NULL,
    chunk_index INT NOT NULL,
    vector_id VARCHAR(100),
    token_count INT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_knowledge_id (knowledge_id),
    FOREIGN KEY (knowledge_id) REFERENCES knowledge_base(id) ON DELETE CASCADE
);

-- 插入示例数据
INSERT IGNORE INTO users (username, email, password_hash, full_name) VALUES 
('admin', 'admin@example.com', '$2a$10$example_hash', '系统管理员'),
('demo', 'demo@example.com', '$2a$10$example_hash', '演示用户');

INSERT IGNORE INTO knowledge_base (title, content, category, tags, created_by) VALUES 
('软件工程师简历模板', '软件工程师简历应该包含以下关键要素...', 'resume_tips', '["软件工程师", "简历模板"]', 1),
('技术技能描述最佳实践', '在简历中描述技术技能时，应该注意以下几点...', 'best_practice', '["技术技能", "最佳实践"]', 1);
EOF
    
    echo "✅ 数据库初始化完成"
}

# 显示启动信息
show_startup_info() {
    echo ""
    echo "🎉 开发环境设置完成！"
    echo ""
    echo "📍 服务地址："
    echo "  - 前端应用: http://localhost:3000"
    echo "  - API网关: http://localhost:8080"
    echo "  - Consul: http://localhost:8500"
    echo "  - MinIO: http://localhost:9001"
    echo "  - Grafana: http://localhost:3001 (admin/admin123)"
    echo "  - Jaeger: http://localhost:16686"
    echo ""
    echo "🚀 启动服务："
    echo "  - 启动前端: cd frontend/web && npm run dev"
    echo "  - 启动后端服务: make dev"
    echo "  - 查看日志: docker-compose -f deployment/docker-compose/docker-compose.yml logs -f"
    echo ""
    echo "⚠️  注意事项："
    echo "  1. 请确保已设置 OPENAI_API_KEY 环境变量"
    echo "  2. 首次启动可能需要较长时间下载镜像"
    echo "  3. 如遇问题请查看日志或重启服务"
    echo ""
}

# 主函数
main() {
    check_prerequisites
    create_configs
    create_env_file
    setup_backend
    setup_frontend
    start_infrastructure
    init_database
    show_startup_info
}

# 执行主函数
main "$@"
