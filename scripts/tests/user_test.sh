#!/bin/bash

run_user_tests() {
    log_info "运行用户测试..."
    
    test_api "获取用户信息" \
        "GET" \
        "/users" \
        "" \
        200 \
        "获取成功"
    
    test_api "更新用户信息" \
        "PUT" \
        "/users" \
        "{\"nickname\":\"测试用户\",\"bio\":\"这是一个测试账号\"}" \
        200 \
        "更新用户信息成功"
    
    test_api "修改密码" \
        "PUT" \
        "/users/password" \
        "{\"oldPassword\":\"$TEST_PASSWORD\",\"newPassword\":\"${TEST_PASSWORD}new\"}" \
        200 \
        "密码修改成功"
} 