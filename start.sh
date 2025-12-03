#!/bin/bash

# VPS Backend Quick Start Script

set -e

echo "======================================"
echo "VPS Backend Management System"
echo "======================================"
echo ""

# 检查 Go 是否安装
if ! command -v go &> /dev/null; then
    echo "❌ Go is not installed. Please install Go 1.21 or higher."
    exit 1
fi

echo "✅ Go version: $(go version)"
echo ""

# 检查 Docker 是否安装
if ! command -v docker &> /dev/null; then
    echo "⚠️  Docker is not installed. You'll need to set up PostgreSQL and Redis manually."
    SKIP_DOCKER=true
else
    echo "✅ Docker is installed"
fi

# 安装依赖
echo "📦 Installing Go dependencies..."
go mod download
go mod tidy

# 启动 Docker 服务
if [ "$SKIP_DOCKER" != "true" ]; then
    echo ""
    echo "🐳 Starting Docker services (PostgreSQL & Redis)..."
    docker-compose up -d
    
    echo "⏳ Waiting for database to be ready..."
    sleep 5
fi

# 运行数据库迁移
echo ""
echo "🔄 Running database migrations..."
go run cmd/server/main.go &
SERVER_PID=$!
sleep 3
kill $SERVER_PID 2>/dev/null || true

# 初始化测试数据
echo ""
read -p "Do you want to seed the database with test data? (y/n) " -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]]; then
    echo "🌱 Seeding database..."
    go run cmd/seed/main.go
fi

# 启动服务器
echo ""
echo "======================================"
echo "🚀 Starting server..."
echo "======================================"
echo ""
echo "Server will start on http://localhost:8080"
echo "API documentation: http://localhost:8080/health"
echo ""
echo "Test credentials:"
echo "  Email: demo@example.com"
echo "  Password: 123456"
echo ""
echo "Admin credentials:"
echo "  Email: admin@example.com"
echo "  Password: admin123456"
echo ""
echo "Press Ctrl+C to stop the server"
echo ""

go run cmd/server/main.go
