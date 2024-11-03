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

run_profile_metadata_tests() {
    log_info "运行个人资料元数据测试..."
    
    # 测试通用信息元数据
    test_api "创建通用信息" \
        "POST" \
        "/profiles" \
        "{\"type\":\"general\",\"title\":\"个人简介\",\"metadata\":{\"display_name\":\"张三\",\"what_you_do\":\"全栈工程师\",\"pronouns\":\"他/他的\",\"about\":\"热爱技术的开发者\"},\"visibility\":\"public\"}" \
        200 \
        "创建成功"

    # 测试项目元数据
    test_api "创建项目经历" \
        "POST" \
        "/profiles" \
        "{\"type\":\"project\",\"title\":\"电商平台\",\"metadata\":{\"client\":\"某科技公司\",\"collaborators\":[\"李四\",\"王五\"],\"attachments\":[{\"type\":\"media\",\"url\":\"https://example.com/demo.mp4\",\"mime_type\":\"video/mp4\",\"size\":1024000,\"name\":\"项目演示.mp4\"}]},\"visibility\":\"public\"}" \
        200 \
        "创建成功"

    # 测试教育经历元数据
    test_api "创建教育经历" \
        "POST" \
        "/profiles" \
        "{\"type\":\"education\",\"title\":\"计算机科学与技术\",\"metadata\":{\"degree\":\"学士\",\"title\":\"本科生\",\"coworkers\":[\"张三\",\"李四\"]},\"visibility\":\"public\"}" \
        200 \
        "创建成功"

    # 测试联系方式元数据
    test_api "创建联系方式" \
        "POST" \
        "/profiles" \
        "{\"type\":\"contact\",\"title\":\"GitHub\",\"metadata\":{\"platform\":\"GitHub\",\"username\":\"zhangsan\",\"email_address\":\"zhangsan@example.com\",\"custom_name\":\"张三的代码库\"},\"visibility\":\"public\"}" \
        200 \
        "创建成功"

    # 测试证书信息元数据
    test_api "创建证书信息" \
        "POST" \
        "/profiles" \
        "{\"type\":\"certification\",\"title\":\"AWS认证\",\"metadata\":{\"issue_date\":\"2023-01-01T00:00:00Z\",\"expiry_date\":\"2026-01-01T00:00:00Z\",\"attachments\":[{\"type\":\"media\",\"url\":\"https://example.com/cert.pdf\",\"mime_type\":\"application/pdf\",\"size\":512000,\"name\":\"证书.pdf\"}]},\"visibility\":\"public\"}" \
        200 \
        "创建成功"

    # 验证元数据查询
    test_api "查询通用信息" \
        "GET" \
        "/profiles?type=general" \
        "" \
        200 \
        "获取成功"

    test_api "查询项目经历" \
        "GET" \
        "/profiles?type=project" \
        "" \
        200 \
        "获取成功"

    # 更新元数据
    test_api "更新教育经历元数据" \
        "PUT" \
        "/profiles/3" \
        "{\"metadata\":{\"degree\":\"硕士\",\"title\":\"研究生\",\"coworkers\":[\"张三\",\"李四\",\"王五\"]}}" \
        200 \
        "更新成功"

    # 验证复杂元数据更新
    test_api "更新项目元数据和附件" \
        "PUT" \
        "/profiles/2" \
        "{\"metadata\":{\"client\":\"某科技公司\",\"collaborators\":[\"李四\",\"王五\",\"赵六\"],\"attachments\":[{\"type\":\"media\",\"url\":\"https://example.com/demo-v2.mp4\",\"mime_type\":\"video/mp4\",\"size\":2048000,\"name\":\"项目演示-新版.mp4\"}]}}" \
        200 \
        "更新成功"
}
