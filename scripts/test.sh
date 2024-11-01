#!/bin/bash

# 获取项目根目录的绝对路径
PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

# 进入项目根目录
cd "${PROJECT_ROOT}"

# 颜色定义
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m'

echo -e "${YELLOW}准备测试环境...${NC}"

# 创建测试配置
cat > .env.test << EOF
DB_HOST=localhost
DB_PORT=5432
DB_NAME=test_db
DB_USER=postgres
DB_PASSWORD=postgres
SERVER_PORT=8080
JWT_SECRET=test-secret
HEALTH_CHECK_INTERVAL=10s
EOF

# 设置测试环境变量
export ENV_FILE=.env.test
export GIN_MODE=test

echo -e "${YELLOW}开始执行单元测试...${NC}\n"

# 执行测试并捕获输出
TEST_OUTPUT=$(go test -v -cover ./internal/... 2>&1)
TEST_RESULT=$?

# 清理测试配置
rm -f .env.test

# 输出测试结果和