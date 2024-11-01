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

echo -e "${YELLOW}开始执行单元测试...${NC}\n"

# 执行测试并捕获输出
TEST_OUTPUT=$(go test -v -cover ./internal/... -coverprofile=coverage.out 2>&1)
TEST_RESULT=$?

# 生成测试报告
if [ -f coverage.out ]; then
    go tool cover -html=coverage.out -o coverage.html
    echo -e "${GREEN}测试报告已生成: coverage.html${NC}"
fi

# 清理测试配置
rm -f .env.test

# 根据测试结果输出相应信息
if [ $TEST_RESULT -eq 0 ]; then
    echo -e "${GREEN}测试通过！${NC}"
else
    echo -e "${RED}测试失败！${NC}"
fi

# 输出测试详细信息
echo -e "\n${YELLOW}测试详细输出：${NC}"
echo "${TEST_OUTPUT}"

exit $TEST_RESULT