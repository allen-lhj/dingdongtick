# Go参数
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod

# 项目参数
BINARY_NAME=ticktick-server
BINARY_PATH=bin/$(BINARY_NAME)
MAIN_PATH=cmd/server/main.go

# 默认目标
.PHONY: all
all: clean deps build

# 安装依赖
.PHONY: deps
deps:
	$(GOMOD) download
	$(GOMOD) tidy

# 构建
.PHONY: build
build:
	mkdir -p bin
	$(GOBUILD) -o $(BINARY_PATH) $(MAIN_PATH)

# 运行
.PHONY: run
run:
	$(GOCMD) run $(MAIN_PATH)

# 运行（生产模式）
.PHONY: run-prod
run-prod:
	GIN_MODE=release $(GOCMD) run $(MAIN_PATH)

# 测试
.PHONY: test
test:
	$(GOTEST) -v ./...

# 测试覆盖率
.PHONY: test-coverage
test-coverage:
	$(GOTEST) -v -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html

# 清理
.PHONY: clean
clean:
	$(GOCLEAN)
	rm -rf bin/
	rm -f coverage.out coverage.html

# 格式化代码
.PHONY: fmt
fmt:
	$(GOCMD) fmt ./...

# 代码检查
.PHONY: vet
vet:
	$(GOCMD) vet ./...

# 安装开发工具
.PHONY: install-tools
install-tools:
	$(GOGET) -u golang.org/x/tools/cmd/goimports
	$(GOGET) -u github.com/golangci/golangci-lint/cmd/golangci-lint

# 代码质量检查
.PHONY: lint
lint:
	golangci-lint run

# 热重载开发（需要安装air）
.PHONY: dev
dev:
	air

# 安装air热重载工具
.PHONY: install-air
install-air:
	$(GOGET) -u github.com/cosmtrek/air

# Docker构建
.PHONY: docker-build
docker-build:
	docker build -t ticktick-backend .

# Docker运行
.PHONY: docker-run
docker-run:
	docker run -p 8080:8080 --env-file .env ticktick-backend

# 数据库相关
.PHONY: db-create
db-create:
	@echo "创建数据库..."
	@psql -h localhost -p 5432 -U postgres -c "CREATE DATABASE ticktick;" || true
	@psql -h localhost -p 5432 -U postgres -c "CREATE USER ticktick_user WITH PASSWORD 'password123';" || true
	@psql -h localhost -p 5432 -U postgres -c "GRANT ALL PRIVILEGES ON DATABASE ticktick TO ticktick_user;" || true

.PHONY: db-drop
db-drop:
	@echo "删除数据库..."
	@psql -h localhost -p 5432 -U postgres -c "DROP DATABASE IF EXISTS ticktick;"
	@psql -h localhost -p 5432 -U postgres -c "DROP USER IF EXISTS ticktick_user;"

.PHONY: db-reset
db-reset: db-drop db-create

# 数据库迁移（需要先启动PostgreSQL）
.PHONY: migrate
migrate:
	$(GOCMD) run $(MAIN_PATH) -migrate

# API 测试
.PHONY: test-api
test-api:
	@echo "测试 API 接口..."
	@echo "1. 健康检查..."
	@curl -s http://localhost:8080/api/v1/health || echo "请确保服务器正在运行"
	@echo ""
	@echo "2. 注册用户..."
	@curl -s -X POST http://localhost:8080/api/v1/auth/register \
		-H "Content-Type: application/json" \
		-d '{"email":"test@example.com","password":"password123","firstName":"Test","lastName":"User"}' \
		-c /tmp/cookies.txt || echo "注册失败或用户已存在"
	@echo ""
	@echo "3. 用户登录..."
	@curl -s -X POST http://localhost:8080/api/v1/auth/login \
		-H "Content-Type: application/json" \
		-d '{"email":"test@example.com","password":"password123"}' \
		-b /tmp/cookies.txt -c /tmp/cookies.txt || echo "登录失败"
	@echo ""
	@echo "4. 获取用户信息..."
	@curl -s -X GET http://localhost:8080/api/v1/profile \
		-b /tmp/cookies.txt || echo "获取用户信息失败"

# 开发环境设置
.PHONY: dev-setup
dev-setup:
	@echo "设置开发环境..."
	@cp .env.example .env || true
	@echo "请编辑 .env 文件设置正确的配置"
	$(MAKE) deps
	$(MAKE) db-create

# 快速启动
.PHONY: quick-start
quick-start: dev-setup run

# 生成API文档（需要安装swag）
.PHONY: docs
docs:
	swag init -g $(MAIN_PATH)

# 帮助
.PHONY: help
help:
	@echo "TickTick Backend 构建工具"
	@echo ""
	@echo "快速开始："
	@echo "  make quick-start   - 快速启动开发环境（推荐）"
	@echo "  make dev-setup     - 设置开发环境"
	@echo ""
	@echo "构建和运行："
	@echo "  make deps          - 安装依赖"
	@echo "  make build         - 构建应用"
	@echo "  make run           - 运行应用（开发模式）"
	@echo "  make run-prod      - 运行应用（生产模式）"
	@echo "  make dev           - 热重载开发"
	@echo ""
	@echo "测试："
	@echo "  make test          - 运行测试"
	@echo "  make test-coverage - 运行测试并生成覆盖率报告"
	@echo "  make test-api      - 测试 API 接口"
	@echo ""
	@echo "数据库："
	@echo "  make db-create     - 创建数据库"
	@echo "  make db-drop       - 删除数据库"
	@echo "  make db-reset      - 重置数据库"
	@echo "  make migrate       - 数据库迁移"
	@echo ""
	@echo "代码质量："
	@echo "  make clean         - 清理构建文件"
	@echo "  make fmt           - 格式化代码"
	@echo "  make vet           - 代码检查"
	@echo "  make lint          - 代码质量检查"
	@echo ""
	@echo "Docker："
	@echo "  make docker-build  - Docker构建"
	@echo "  make docker-run    - Docker运行"
	@echo ""
	@echo "其他："
	@echo "  make docs          - 生成API文档"
	@echo "  make help          - 显示帮助信息"
