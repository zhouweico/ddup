#!/bin/bash

# 检查是否提供了提交信息
if [ -z "$1" ]; then
    echo "请提供提交信息"
    echo "使用方法: ./scripts/commit.sh '提交信息'"
    exit 1
fi

# 生成 swagger 文档
echo "正在生成 Swagger 文档..."
swag init -g cmd/api/main.go -o docs

# 检查 swag 命令是否执行成功
if [ $? -ne 0 ]; then
    echo "Swagger 文档生成失败"
    exit 1
fi

# 添加所有更改到暂存区
git add .

# 提交更改
git commit -m "$1"

# 检查提交是否成功
if [ $? -ne 0 ]; then
    echo "Git 提交失败"
    exit 1
fi

# 推送到远程仓库
git push

# 检查推送是否成功
if [ $? -ne 0 ]; then
    echo "推送到远程仓库失败"
    exit 1
fi

echo "完成！"
echo "- Swagger 文档已更新"
echo "- 更改已提交: $1"
echo "- 更改已推送到远程仓库"