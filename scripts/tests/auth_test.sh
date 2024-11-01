#!/bin/bash

# 注册测试函数
run_register_test() {
    log_info "运行注册测试..."
    
    test_api "用户注册" \
        "POST" \
        "/auth/register" \
        "{\"username\":\"$TEST_USER\",\"password\":\"$TEST_PASSWORD\",\"email\":\"$TEST_USER@example.com\"}" \
        200 \
        "注册成功" || true  # 即使注册失败也继续执行
}

# 登录测试函数
run_auth_tests() {
    log_info "运行登录测试..."
    
    # 登录测试
    local login_response=$(curl -s -X POST \
        -H "Content-Type: application/json" \
        -d "{\"username\":\"$TEST_USER\",\"password\":\"$TEST_PASSWORD\"}" \
        "$API_URL/auth/login")
    
    # 提取 token
    export TOKEN=$(echo $login_response | jq -r '.data.token')
    
    # 验证登录是否成功
    if [ -z "$TOKEN" ] || [ "$TOKEN" = "null" ]; then
        log_error "登录失败"
        return 1
    fi
    
    log_success "登录成功"
    return 0
}

# 登出测试函数
run_logout_test() {
    log_info "运行登出测试..."
    
    test_api "用户登出" \
        "POST" \
        "/auth/logout" \
        "" \
        200 \
        "退出成功"
} 