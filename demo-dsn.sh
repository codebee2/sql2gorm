#!/bin/bash

# SQL2GORM DSN功能演示脚本
# 此脚本演示如何使用DSN连接数据库生成GORM模型

echo "🚀 SQL2GORM DSN功能演示"
echo "================================"

# 检查是否已构建
if [ ! -f "./bin/sql2gorm" ]; then
    echo "⚠️  工具未构建，正在构建..."
    make build
fi

echo ""
echo "📋 可用命令示例："
echo ""

echo "1️⃣  使用DSN配置文件："
echo "   ./bin/sql2gorm --config config-dsn.yaml"
echo ""

echo "2️⃣  使用命令行DSN参数："
echo "   ./bin/sql2gorm \\"
echo "     --dsn \"root:password@tcp(localhost:3306)/testdb?charset=utf8mb4&parseTime=True&loc=Local\" \\"
echo "     --output models/user.go \\"
echo "     --table users \\"
echo "     --package models"
echo ""

echo "3️⃣  使用Makefile DSN命令："
echo "   make run-dsn"
echo ""

echo "4️⃣  查看帮助信息："
echo "   ./bin/sql2gorm --help"
echo ""

echo "📝 注意事项："
echo "   - 请确保数据库服务正在运行"
echo "   - 请修改DSN中的用户名、密码、数据库名等信息"
echo "   - 确保数据库中存在指定的表"
echo ""

echo "🔧 配置DSN："
echo "   格式：username:password@tcp(host:port)/database?charset=utf8mb4&parseTime=True&loc=Local"
echo "   示例：root:mypassword@tcp(127.0.0.1:3306)/mydb?charset=utf8mb4&parseTime=True&loc=Local"
echo ""

echo "✅ 演示脚本完成！" 