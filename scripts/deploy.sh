#!/bin/bash

set -e

echo "=== 强制学习系统部署脚本 ==="

# 检查 Docker
if ! command -v docker &> /dev/null; then
    echo "错误: Docker 未安装"
    exit 1
fi

if ! command -v docker-compose &> /dev/null; then
    echo "错误: Docker Compose 未安装"
    exit 1
fi

# 启动服务
echo "1. 启动 PostgreSQL 和 Redis..."
docker-compose up -d postgres redis

# 等待 PostgreSQL 就绪
echo "2. 等待 PostgreSQL 就绪..."
sleep 5

# 检查 PostgreSQL 健康状态
for i in {1..30}; do
    if docker-compose exec -T postgres pg_isready -U postgres &> /dev/null; then
        echo "   PostgreSQL 已就绪"
        break
    fi
    echo "   等待 PostgreSQL... ($i/30)"
    sleep 2
done

# 构建并启动后端
echo "3. 构建并启动后端服务..."
docker-compose up -d --build backend

# 检查后端健康状态
echo "4. 检查后端服务状态..."
sleep 5

# 显示服务状态
echo ""
echo "=== 服务状态 ==="
docker-compose ps

echo ""
echo "=== 访问信息 ==="
echo "后端 API: http://localhost:8080"
echo "健康检查: http://localhost:8080/health"

echo ""
echo "部署完成!"
