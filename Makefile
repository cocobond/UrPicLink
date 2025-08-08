.PHONY: build run test clean docker-build docker-run help

# 默认目标
.DEFAULT_GOAL := help

# 应用名称
APP_NAME := urpicbed

# 构建应用
build:
	@echo "构建 $(APP_NAME)..."
	go build -o $(APP_NAME) .

# 运行应用
run: build
	@echo "运行 $(APP_NAME)..."
	./$(APP_NAME)

# 测试
test:
	@echo "运行测试..."
	go test ./...

# 清理构建文件
clean:
	@echo "清理构建文件..."
	rm -f $(APP_NAME)
	rm -f $(APP_NAME).exe

# 下载依赖
deps:
	@echo "下载依赖..."
	go mod download
	go mod tidy

# Docker构建
docker-build:
	@echo "构建Docker镜像..."
	docker build -t $(APP_NAME) .

# Docker运行
docker-run: docker-build
	@echo "运行Docker容器..."
	docker run -p 8080:8080 -v $(PWD)/config:/root/config $(APP_NAME)

# Docker Compose启动
docker-compose-up:
	@echo "启动Docker Compose服务..."
	docker-compose up -d

# Docker Compose停止
docker-compose-down:
	@echo "停止Docker Compose服务..."
	docker-compose down

# Docker Compose重启
docker-compose-restart:
	@echo "重启Docker Compose服务..."
	docker-compose restart

# 查看日志
logs:
	@echo "查看服务日志..."
	docker-compose logs -f

# 格式化代码
fmt:
	@echo "格式化代码..."
	go fmt ./...

# 代码检查
lint:
	@echo "代码检查..."
	golangci-lint run

# 帮助信息
help:
	@echo "可用的命令:"
	@echo "  build              - 构建应用"
	@echo "  run                - 构建并运行应用"
	@echo "  test               - 运行测试"
	@echo "  clean              - 清理构建文件"
	@echo "  deps               - 下载依赖"
	@echo "  docker-build       - 构建Docker镜像"
	@echo "  docker-run         - 构建并运行Docker容器"
	@echo "  docker-compose-up  - 启动Docker Compose服务"
	@echo "  docker-compose-down - 停止Docker Compose服务"
	@echo "  docker-compose-restart - 重启Docker Compose服务"
	@echo "  logs               - 查看服务日志"
	@echo "  fmt                - 格式化代码"
	@echo "  lint               - 代码检查"
	@echo "  help               - 显示此帮助信息" 