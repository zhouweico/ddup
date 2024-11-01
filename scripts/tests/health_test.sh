#!/bin/bash

run_health_tests() {
    log_info "运行健康检查测试..."
    
    # 注意：健康检查使用根路径
    local health_response=$(curl -s -X GET "http://localhost:8080/health")
    
    if ! echo "$health_response" | jq . >/dev/null 2>&1; then
        log_error "无效的 JSON 响应"
        return 1
    fi
    
    local status=$(echo "$health_response" | jq -r '.code')
    local message=$(echo "$health_response" | jq -r '.message')
    
    if [ "$status" -eq 200 ] && [ "$message" = "服务正常" ]; then
        log_success "健康检查测试通过"
        PASSED_TESTS=$((PASSED_TESTS + 1))
        return 0
    else
        log_error "健康检查测试失败"
        return 1
    fi
    
    TOTAL_TESTS=$((TOTAL_TESTS + 1))
} 