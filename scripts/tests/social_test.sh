#!/bin/bash

run_social_tests() {
    log_info "运行社交媒体测试..."
    
    # 创建社交媒体账号
    test_api "创建社交媒体账号" \
        "POST" \
        "/users/socials" \
        "{\"platform\":\"github\",\"username\":\"testuser\",\"url\":\"https://github.com/testuser\"}" \
        200 \
        "创建成功"
    
    # 获取社交媒体账号列表
    test_api "获取社交媒体账号列表" \
        "GET" \
        "/users/socials" \
        "" \
        200 \
        "获取成功"
    
    # 更新社交媒体账号
    test_api "更新社交媒体账号" \
        "PUT" \
        "/users/socials/1" \
        "{\"platform\":\"github\",\"username\":\"testuser-updated\"}" \
        200 \
        "更新成功"
    
    # 删除社交媒体账号
    test_api "删除社交媒体账号" \
        "DELETE" \
        "/users/socials/1" \
        "" \
        200 \
        "删除成功"
} 