#!/bin/bash

# API 配置
API_URL="http://localhost:8080/api/v1"
TEST_USER="testuser"
TEST_PASSWORD="password123"

# 测试结果统计
TOTAL_TESTS=0
PASSED_TESTS=0

# 颜色定义
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m'

# 日志函数
log_info() { echo -e "${YELLOW}[INFO] $1${NC}"; }
log_success() { echo -e "${GREEN}[SUCCESS] $1${NC}"; }
log_error() { echo -e "${RED}[ERROR] $1${NC}"; }

# API 测试函数
test_api() {
    local name=$1
    local method=$2
    local endpoint=$3
    local data=$4
    local expected_status=$5
    local expected_message=$6
    
    TOTAL_TESTS=$((TOTAL_TESTS + 1))
    
    log_info "测试: $name"
    
    # 构建 curl 命令
    local curl_cmd="curl -s -X $method"
    [[ -n "$data" ]] && curl_cmd="$curl_cmd -H 'Content-Type: application/json' -d '$data'"
    [[ -n "$TOKEN" ]] && curl_cmd="$curl_cmd -H 'Authorization: $TOKEN'"
    curl_cmd="$curl_cmd $API_URL$endpoint"
    
    # 执行请求
    local response=$(eval $curl_cmd)
    verify_response "$response" "$expected_status" "$expected_message"
}

# 响应验证函数
verify_response() {
    local response=$1
    local expected_status=$2
    local expected_message=$3
    
    if ! echo "$response" | jq . >/dev/null 2>&1; then
        log_error "无效的 JSON 响应"
        return 1
    fi
    
    local status=$(echo "$response" | jq -r '.code // empty')
    local message=$(echo "$response" | jq -r '.message // empty')
    
    if [ -z "$status" ]; then
        log_error "响应中没有找到 code 字段"
        return 1
    fi
    
    if ! [[ "$status" =~ ^[0-9]+$ ]]; then
        log_error "无效的状态码: $status"
        return 1
    fi
    
    if [ "$status" -eq "$expected_status" ] && [ "$message" = "$expected_message" ]; then
        log_success "测试通过"
        PASSED_TESTS=$((PASSED_TESTS + 1))
        return 0
    else
        log_error "测试失败： $message"
        return 1
    fi
}

# 打印测试结果
print_test_results() {
    echo -e "\n${YELLOW}测试完成${NC}"
    echo -e "总计: ${TOTAL_TESTS}"
    echo -e "通过: ${GREEN}${PASSED_TESTS}${NC}"
    echo -e "失败: ${RED}$((TOTAL_TESTS - PASSED_TESTS))${NC}"
} 