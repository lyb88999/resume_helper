# ResumeOptim 多服务项目 Makefile

# --- 变量定义 ---
# Go 命令
GOPATH:=$(shell go env GOPATH)
GO_CMD=go

# Git 版本号
VERSION?=$(shell git describe --tags --always --dirty)

# 服务列表 (未来新增服务时在此处添加即可)
SERVICES := user-service file-service parser-service ai-service

# Proto 文件路径
# 只编译项目自身的API定义，不包括第三方依赖
#API_PROTO_FILES := $(shell find api -name *.proto)
API_PROTO_FILES := $(shell find api -name "*.proto") $(shell find backend/services -path "*/api/*.proto")

# 只编译每个服务的conf.proto和shared/proto下的文件
INTERNAL_PROTO_FILES := $(shell find backend/services/*/internal/conf -name *.proto) $(shell find backend/shared/proto -name *.proto)


# --- 开发工具与环境 ---

.PHONY: help
# 显示帮助信息
help:
	@echo ''
	@echo '用法:'
	@echo ' make [target]'
	@echo ''
	@echo '可用目标:'
	@awk '/^[a-zA-Z\-\_0-9]+:/ { \
	helpMessage = match(lastLine, /^# (.*)/); \
		if (helpMessage) { \
			helpCommand = substr($$1, 0, index($$1, ":")); \
			helpMessage = substr(lastLine, RSTART + 2, RLENGTH); \
			printf "\033[36m%-22s\033[0m %s\n", helpCommand,helpMessage; \
		} \
	} \
	{ lastLine = $$0 }' $(MAKEFILE_LIST)

.PHONY: init
# 安装所有必要的Go代码生成工具
init:
	$(GO_CMD) install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	$(GO_CMD) install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	$(GO_CMD) install github.com/go-kratos/kratos/cmd/protoc-gen-go-http/v2@latest
	$(GO_CMD) install github.com/go-kratos/kratos/cmd/protoc-gen-go-errors/v2@latest
	$(GO_CMD) install github.com/google/wire/cmd/wire@latest

.PHONY: setup
# 一键完成新环境的初始化设置
setup: init all
	@echo "✅ 环境设置完成。请配置您的 .env 文件，然后运行 'make dev' 来启动所有服务。"


# --- 代码生成 ---

.PHONY: proto
# 生成所有Protobuf相关代码
proto: api config

.PHONY: api
# 生成项目根目录下的 /api/**/*.proto 文件
api:
	@echo "正在生成 API proto 文件..."
	protoc --proto_path=. \
	       --proto_path=./third_party \
	       --go_out=paths=source_relative:. \
	       --go-http_out=paths=source_relative:. \
	       --go-grpc_out=paths=source_relative:. \
	       --go-errors_out=paths=source_relative:. \
	       $(API_PROTO_FILES)

.PHONY: config
# 生成所有内部和共享的proto文件
config:
	@echo "正在生成内部和共享 proto 文件..."
	protoc --proto_path=. \
	       --proto_path=./third_party \
	       --go_out=paths=source_relative:. \
	       $(INTERNAL_PROTO_FILES)

.PHONY: wire
# 为所有服务生成Wire依赖注入代码
wire:
	@echo "正在为所有服务生成 Wire 依赖注入代码..."
	@for service in $(SERVICES); do \
		echo "  -> 正在处理服务：$$service..."; \
		cd backend/services/$$service/cmd/*-service && $(GO_CMD) run github.com/google/wire/cmd/wire; \
		cd ../../../../..; \
	done

.PHONY: generate
# 整理Go模块并运行go generate
generate:
	@echo "正在整理 Go 模块..."
	$(GO_CMD) mod tidy
	@echo "正在运行 go generate..."
	$(GO_CMD) generate ./...

.PHONY: all
# 一键生成所有代码
all: proto generate wire


# --- 构建与运行 ---

.PHONY: build
# 构建所有微服务到./bin目录
build:
	@echo "🔨 正在构建所有服务..."
	@mkdir -p bin/
	@for service in $(SERVICES); do \
		echo "  -> 正在构建：$$service..."; \
		$(GO_CMD) build -ldflags "-X main.Version=$(VERSION)" -o ./bin/$$service ./backend/services/$$service/cmd/*-service; \
	done
	@echo "✅ 所有服务构建完成。"

# 构建指定服务, e.g., make build-user-service
build-%:
	@echo "🔨 正在构建服务: $*..."
	@mkdir -p bin/
	$(GO_CMD) build -ldflags "-X main.Version=$(VERSION)" -o ./bin/$* ./backend/services/$*/cmd/*-service

# 本地运行指定服务, e.g., make run-user-service
.PHONY: run-%
run-%:
	@echo "🚀 正在运行服务: $*..."
	@$(GO_CMD) run ./backend/services/$*/cmd/*-service -conf ./configs/$*.yaml

.PHONY: dev
# (推荐) 在本地并发运行所有后端服务
dev: build
	@echo "🚀 正在并发启动所有后端服务..."
	@for service in $(SERVICES); do \
		echo "  -> 正在后台启动 $$service..."; \
		./bin/$$service -conf ./configs/$$service.yaml & \
	done
	@echo "✅ 所有服务已在后台启动。使用 'make stop-dev' 来停止它们。"
	@wait

.PHONY: stop-dev
# 停止由 'make dev' 启动的所有服务
stop-dev:
	@echo "🛑 正在停止所有后台开发服务..."
	@pkill -f "./bin/user-service" || true
	@pkill -f "./bin/file-service" || true
	@pkill -f "./bin/parser-service" || true
	@pkill -f "./bin/ai-service" || true
	@echo "✅ 所有服务已停止。"


# --- 清理 ---

.PHONY: clean
# 清理所有生成的文件和构建产物
clean:
	@echo "🧹 正在清理项目..."
	rm -rf bin
	find . -name "*.pb.go" -type f -delete
	find . -name "*.pb.http.go" -type f -delete
	find . -name "*.pb.grpc.go" -type f -delete
	find . -name "*.pb.errors.go" -type f -delete
	find . -name "wire_gen.go" -type f -delete
	find . -name "openapi.yaml" -type f -delete


# --- Docker ---

.PHONY: docker-up
# 使用docker-compose启动开发环境
docker-up:
	docker-compose -f docker-compose.dev.yml up -d

.PHONY: docker-down
# 关闭docker-compose开发环境
docker-down:
	docker-compose -f docker-compose.dev.yml down

.PHONY: docker-logs
# 查看所有服务的docker日志
docker-logs:
	docker-compose -f docker-compose.dev.yml logs -f



.DEFAULT_GOAL := help