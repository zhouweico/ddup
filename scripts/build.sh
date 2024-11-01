#!/bin/bash

# 设置版本号
VERSION=$(git describe --tags --always --dirty)
IMAGE_NAME="ddup-apis"

# 生成 swagger 文档
echo "生成 Swagger 文档..."
swag init -g cmd/api/main.go -o docs

# 清理旧镜像
echo "清理旧镜像..."
docker rmi ${IMAGE_NAME}:latest >/dev/null 2>&1 || true
docker rmi ${IMAGE_NAME}:${VERSION} >/dev/null 2>&1 || true

# 构建 Docker 镜像
echo "构建 Docker 镜像..."
docker build -t ${IMAGE_NAME}:${VERSION} .
docker tag ${IMAGE_NAME}:${VERSION} ${IMAGE_NAME}:latest

echo "构建完成！"
echo "镜像标签: ${IMAGE_NAME}:${VERSION}"