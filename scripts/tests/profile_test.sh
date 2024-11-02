#!/bin/bash

run_profile_tests() {
    log_info "运行个人资料测试..."
    
    # 创建教育经历
    test_api "创建教育经历" \
        "POST" \
        "/profiles" \
        "{\"type\":\"education\",\"title\":\"测试大学\",\"organization\":\"测试大学\",\"description\":\"计算机科学与技术\",\"start_date\":\"2020-09-01T00:00:00Z\",\"end_date\":\"2024-06-30T00:00:00Z\",\"location\":\"北京\",\"visibility\":\"public\"}" \
        200 \
        "创建成功"
    
    # 创建工作经历
    test_api "创建工作经历" \
        "POST" \
        "/profiles" \
        "{\"type\":\"work\",\"title\":\"软件工程师\",\"organization\":\"测试公司\",\"description\":\"全栈开发\",\"start_date\":\"2024-07-01T00:00:00Z\",\"location\":\"上海\",\"visibility\":\"public\"}" \
        200 \
        "创建成功"
    
    # 创建项目经历
    test_api "创建项目经历" \
        "POST" \
        "/profiles" \
        "{\"type\":\"project\",\"title\":\"开源项目\",\"organization\":\"GitHub\",\"description\":\"一个开源项目\",\"start_date\":\"2023-01-01T00:00:00Z\",\"end_date\":\"2023-12-31T00:00:00Z\",\"url\":\"https://github.com/test/project\",\"visibility\":\"public\"}" \
        200 \
        "创建成功"
    
    # 获取不同类型的个人资料列表
    test_api "获取教育经历" \
        "GET" \
        "/profiles?type=education" \
        "" \
        200 \
        "获取成功"
    
    test_api "获取工作经历" \
        "GET" \
        "/profiles?type=work" \
        "" \
        200 \
        "获取成功"
    
    test_api "获取项目经历" \
        "GET" \
        "/profiles?type=project" \
        "" \
        200 \
        "获取成功"
    
    # 更新个人资料
    test_api "更新教育经历" \
        "PUT" \
        "/profiles/1" \
        "{\"title\":\"测试大学(已更新)\",\"description\":\"软件工程\"}" \
        200 \
        "更新成功"
    
    # 更新显示顺序
    test_api "更新显示顺序" \
        "PUT" \
        "/profiles/order" \
        "{\"items\":[{\"id\":3,\"order\":1},{\"id\":2,\"order\":2},{\"id\":1,\"order\":3}]}" \
        200 \
        "更新显示顺序成功"
    
    # 删除个人资料
    test_api "删除个人资料" \
        "DELETE" \
        "/profiles/1" \
        "" \
        200 \
        "删除成功"
} 