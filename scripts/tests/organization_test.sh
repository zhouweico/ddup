#!/bin/bash

run_organization_tests() {
    log_info "运行组织测试..."
    
    # 创建组织
    test_api "创建组织" \
        "POST" \
        "/organizations" \
        "{\"name\":\"测试组织\",\"description\":\"这是一个测试组织\"}" \
        200 \
        "创建组织成功"
    
    # 获取组织列表
    test_api "获取组织列表" \
        "GET" \
        "/organizations" \
        "" \
        200 \
        "获取组织列表成功"
    
    # 更新组织信息
    test_api "更新组织信息" \
        "PUT" \
        "/organizations/1" \
        "{\"name\":\"更新后的组织\",\"description\":\"更新后的描述\"}" \
        200 \
        "更新组织信息成功"
    
    # 获取组织成员列表
    test_api "获取组织成员列表" \
        "GET" \
        "/organizations/1/members" \
        "" \
        200 \
        "获取成员列表成功"
    
    # 添加组织成员
    test_api "添加组织成员" \
        "POST" \
        "/organizations/1/members" \
        "{\"userID\":2,\"role\":\"member\"}" \
        200 \
        "添加成员成功"
    
    # 删除组织
    test_api "删除组织" \
        "DELETE" \
        "/organizations/1" \
        "" \
        200 \
        "删除组织成功"
} 