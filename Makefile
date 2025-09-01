.PHONY: build test run clean help

# 默认目标
.DEFAULT_GOAL := help

# 构建可执行文件
build:
	@echo "🔨 构建SQL2GORM工具..."
	go build -o bin/sql2gorm main.go
	@echo "✅ 构建完成: bin/sql2gorm"

# 运行测试
test:
	@echo "🧪 运行测试..."
	go test ./...
	@echo "✅ 测试完成"

# 运行工具（使用配置文件）
run:
	@echo "🚀 运行SQL2GORM工具..."
	go run main.go



# 安装依赖
deps:
	@echo "📦 安装依赖..."
	go mod tidy
	@echo "✅ 依赖安装完成"

# 清理构建文件
clean:
	@echo "🧹 清理构建文件..."
	rm -rf bin/
	@echo "✅ 清理完成"

# 显示帮助信息
help:
	@echo "SQL2GORM - SQL表结构转GORM模型工具"
	@echo ""
	@echo "可用命令:"
	@echo "  make build      - 构建可执行文件"
	@echo "  make test       - 运行测试"
	@echo "  make run        - 运行工具（使用配置文件）"
	@echo "  make run-cli    - 运行工具（命令行模式）"
	@echo "  make run-dsn    - 运行工具（DSN模式）"
	@echo "  make deps       - 安装依赖"
	@echo "  make clean      - 清理构建文件"
	@echo "  make help       - 显示此帮助信息"
	@echo ""
	@echo "示例:"
	@echo "  make build      # 构建工具"
	@echo "  make run-cli    # 使用命令行参数运行"
	@echo "  make run-dsn    # 使用DSN连接数据库运行"
	@echo "  ./bin/sql2gorm --help  # 查看帮助" 