#!/bin/bash

# 获取项目根目录的绝对路径
PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
ENV_FILE="${PROJECT_ROOT}/.env"

# 检查 .env 文件是否存在
if [ ! -f "${ENV_FILE}" ]; then
    echo "错误: .env 文件不存在于 ${ENV_FILE}"
    echo "请先复制 .env.example 并配置环境变量"
    exit 1
fi

# 进入项目根目录
cd "${PROJECT_ROOT}"

# 先构建镜像
echo "构建镜像..."
./scripts/build.sh
if [ $? -ne 0 ]; then
    echo "错误: 镜像构建失败"
    exit 1
fi

# 启动服务
echo "启动服务..."
docker compose up -d

# 检查服务状态
echo "检查服务状态..."
if ! docker compose ps --format json | grep -q "running"; then
    echo "错误: 部分服务启动失败"
    docker compose logs
    exit 1
fi

echo "等待数据库就绪..."
for i in {1..30}; do
    if docker compose exec postgres pg_isready -U ddup >/dev/null 2>&1; then
        echo "数据库已就绪"
        break
    fi
    echo -n "."
    sleep 1
done

if [ $i -eq 30 ]; then
    echo "错误: 数据库启动超时"
    exit 1
fi

echo "部署完成！"
echo "服务访问地址: http://localhost:8080"
echo ""
echo "可用命令："
echo "查看日志: docker compose logs -f"
echo "停止服务: docker compose down"
echo "重启服务: docker compose restart"