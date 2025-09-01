#!/bin/bash

echo "🚀 SQL2GORM 工具演示"
echo "======================"

# 检查是否已构建
if [ ! -f "bin/sql2gorm" ]; then
    echo "📦 构建工具..."
    make build
fi

echo ""
echo "1️⃣ 查看帮助信息"
echo "----------------"
./bin/sql2gorm --help

echo ""
echo "2️⃣ 使用配置文件运行"
echo "------------------"
./bin/sql2gorm

echo ""
echo "3️⃣ 使用命令行参数运行"
echo "--------------------"
./bin/sql2gorm \
    --sql example.sql \
    --output models/demo_output.go \
    --table ny_internal_miniprogram \
    --package models

echo ""
echo "4️⃣ 查看生成的代码"
echo "----------------"
if [ -f "models/demo_output.go" ]; then
    echo "✅ 成功生成文件: models/demo_output.go"
    echo ""
    echo "文件内容预览:"
    echo "============="
    head -20 models/demo_output.go
    echo "..."
else
    echo "❌ 文件生成失败"
fi

echo ""
echo "5️⃣ 运行测试"
echo "----------"
make test

echo ""
echo "🎉 演示完成！"
echo ""
echo "📁 生成的文件:"
ls -la models/ 2>/dev/null || echo "models目录不存在"
echo ""
echo "💡 提示: 使用 'make help' 查看所有可用命令" 