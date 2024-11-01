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

# API基础URL
API_URL="http://localhost:8080"
TEST_USER="testuser"
TEST_PASSWORD="password123"
TOKEN=""

# 测试结果统计
TOTAL_TESTS=0
PASSED_TESTS=0

# 执行API请求并检查响应
test_api() {
    local name=$1
    local method=$2
    local endpoint=$3
    local data=$4
    local expected_status=$5
    local expected_message=$6
    
    TOTAL_TESTS=$((TOTAL_TESTS + 1))
    
    echo -e "\n${YELLOW}测试: $name${NC}"
    
    # 构建curl命令
    local curl_cmd="curl -s -X $method"
    if [ ! -z "$data" ]; then
        curl_cmd="$curl_cmd -H 'Content-Type: application/json' -d '$data'"
    fi
    if [ ! -z "$TOKEN" ]; then
        curl_cmd="$curl_cmd -H 'Authorization: $TOKEN'"
    fi
    curl_cmd="$curl_cmd $API_URL$endpoint"
    
    # 执行请求并保存响应
    local response=$(eval $curl_cmd)
    
    # 调试输出
    echo "响应内容: $response"
    
    # 检查响应是否为有效的 JSON
    if ! echo "$response" | jq . >/dev/null 2>&1; then
        echo -e "${RED}✗ 失败: 无效的 JSON 响应${NC}"
        echo "原始响应: $response"
        return 1
    fi
    
    # 解析响应
    local status=$(echo "$response" | jq -r '.code // empty')
    local message=$(echo "$response" | jq -r '.message // empty')
    
    # 检查是否成功解析
    if [ -z "$status" ] || [ -z "$message" ]; then
        echo -e "${RED}✗ 失败: 无法解析响应${NC}"
        echo "状态码: $status"
        echo "消息: $message"
        return 1
    fi
    
    # 检查响应
    if [ "$status" -eq "$expected_status" ] && [ "$message" = "$expected_message" ]; then
        echo -e "${GREEN}✓ 通过${NC}"
        PASSED_TESTS=$((PASSED_TESTS + 1))
    else
        echo -e "${RED}✗ 失败${NC}"
        echo "期望状态码: $expected_status, 实际: $status"
        echo "期望消息: $expected_message, 实际: $message"
    fi
}

# 启动测试
echo -e "${YELLOW}开始API测试...${NC}"

# 注册测试
test_api "注册新用户" "POST" "/api/v1/register" \
    "{\"username\":\"$TEST_USER\",\"password\":\"$TEST_PASSWORD\"}" \
    200 "注册成功"

# 重复注册测试
test_api "重复注册" "POST" "/api/v1/register" \
    "{\"username\":\"$TEST_USER\",\"password\":\"$TEST_PASSWORD\"}" \
    400 "用户名已存在"

test_api "登录" "POST" "/api/v1/login" \
    "{\"username\":\"$TEST_USER\",\"password\":\"$TEST_PASSWORD\"}" \
    200 "登录成功"

# 登录测试
response=$(curl -s -X POST "$API_URL/api/v1/login" \
    -H 'Content-Type: application/json' \
    -d "{\"username\":\"$TEST_USER\",\"password\":\"$TEST_PASSWORD\"}")

# 检查登录响应
if echo "$response" | jq . >/dev/null 2>&1; then
    TOKEN=$(echo "$response" | jq -r '.data.token // empty')
    if [ -z "$TOKEN" ]; then
        echo -e "${RED}警告: 无法获取登录令牌${NC}"
    fi
else
    echo -e "${RED}错误: 登录响应不是有效的 JSON${NC}"
    echo "响应内容: $response"
fi

# 更新用户信息测试
test_api "更新用户信息" "PUT" "/api/v1/user" \
    "{\"nickname\":\"测试用户\",\"email\":\"test@example.com\",\"mobile\":\"18888888888\",\"avatar\":\"https://example.com/avatar.png\",\"bio\":\"测试用户\",\"location\":\"测试地点\",\"website\":\"https://example.com\"}" \
    200 "更新用户信息成功"

# 获取用户信息测试
test_api "获取用户信息" "GET" "/api/v1/user" "" \
    200 "获取成功"

# 修改密码测试
test_api "修改密码" "PUT" "/api/v1/user/password" \
    "{\"oldPassword\":\"$TEST_PASSWORD\",\"newPassword\":\"newpassword123\"}" \
    200 "密码修改成功"

# 创建社交媒体账号测试
test_api "创建社交媒体账号" "POST" "/api/v1/social" \
    "{\"platform\":\"github\",\"username\":\"testuser\",\"url\":\"https://github.com/testuser\",\"description\":\"我的 GitHub 账号\"}" \
    200 "创建成功"

# 创建另一个社交媒体账号
test_api "创建第二个社交媒体账号" "POST" "/api/v1/social" \
    "{\"platform\":\"twitter\",\"username\":\"testuser\",\"url\":\"https://twitter.com/testuser\",\"description\":\"我的 Twitter 账号\"}" \
    200 "创建成功"

# 获取用户的社交媒体账号列表
test_api "获取社交媒体账号列表" "GET" "/api/v1/social" "" \
    200 "获取成功"

# 保存第一个社交媒体账号的ID用于后续测试
response=$(curl -s -X GET "$API_URL/api/v1/social" \
    -H "Authorization: $TOKEN")
SOCIAL__ID=$(echo "$response" | jq -r '.data[0].id // empty')

if [ ! -z "$SOCIAL__ID" ]; then
    # 更新社交媒体账号测试
    test_api "更新社交媒体账号" "PUT" "/api/v1/social/$SOCIAL__ID" \
        "{\"platform\":\"github\",\"username\":\"testuser_updated\",\"url\":\"https://github.com/testuser_updated\",\"description\":\"更新后的 GitHub 账号\"}" \
        200 "更新成功"

    # 删除社交媒体账号测试
    test_api "删除社交媒体账号" "DELETE" "/api/v1/social/$SOCIAL__ID" "" \
        200 "删除成功"
else
    echo -e "${RED}警告: 无法获取社交媒体账号 ID${NC}"
fi

# 测试无效的社交媒体账号ID
test_api "测试无效的社交媒体账号ID" "GET" "/api/v1/social/999999" "" \
    404 "社交媒体账号不存在"

# 测试创建无效的社交媒体平台
test_api "创建无效的社交媒体平台" "POST" "/api/v1/social" \
    "{\"platform\":\"\",\"username\":\"testuser\"}" \
    400 "无效的请求参数"

# 登出测试
test_api "登出" "POST" "/api/v1/logout" "" \
    200 "退出成功"

# 输出测试结果统计
echo -e "\n${YELLOW}测试完成${NC}"
echo -e "总计: $TOTAL_TESTS"
echo -e "通过: ${GREEN}$PASSED_TESTS${NC}"
echo -e "失败: ${RED}$((TOTAL_TESTS - PASSED_TESTS))${NC}"

# 设置退出状态
if [ $PASSED_TESTS -eq $TOTAL_TESTS ]; then
    exit 0
else
    exit 1
fi 