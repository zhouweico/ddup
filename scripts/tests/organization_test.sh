#!/bin/bash

run_organization_tests() {
    log_info "运行组织测试..."
    
    # 使用时间戳生成唯一的组织名称
    local timestamp=$(date +%s)
    local org_name="test-org"
    local display_name="测试组织"
    
    # 创建组织
    test_api "创建组织" \
        "POST" \
        "/organizations" \
        "{\"name\":\"test-org\",\"display_name\":\"测试组织\",\"description\":\"这是一个测试组织\"}" \
        200 \
        "创建组织成功"
    
    # 测试组织名称唯一性
    test_api "创建重复组织" \
        "POST" \
        "/organizations" \
        "{\"name\":\"test-org\",\"display_name\":\"测试组织\",\"description\":\"这是一个重复的组织\"}" \
        400 \
        "组织名称已存在"
    
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
        "/organizations/test-org" \
        "{\"display_name\":\"测试组织 已更新\",\"description\":\"更新后的描述\"}" \
        200 \
        "更新组织信息成功"
    
    # 获取组织成员列表
    test_api "获取组织成员列表" \
        "GET" \
        "/organizations/test-org/members" \
        "" \
        200 \
        "获取成员列表成功"
    
    # 添加组织成员
    test_api "添加组织成员" \
        "POST" \
        "/organizations/test-org/members" \
        "{\"username\":\"testuser\",\"role\":\"member\"}" \
        200 \
        "添加成员成功"
    
    # 测试加入组织
    test_api "加入组织" \
        "POST" \
        "/organizations/test-org/join" \
        "" \
        200 \
        "加入组织成功"
    
    # 更新成员信息成功
    test_api "更新成员信息成功" \
        "PUT" \
        "/organizations/test-org/members/testuser" \
        "{\"role\":\"admin\"}" \
        200 \
        "更新成员信息成功"
    
    # 移除组织成员
    test_api "移除组织成员" \
        "DELETE" \
        "/organizations/test-org/members/testuser" \
        "" \
        200 \
        "移除成员成功"
    
    # 删除组织
    test_api "删除组织" \
        "DELETE" \
        "/organizations/test-org" \
        "" \
        200 \
        "删除组织成功"
} 