#!/bin/bash

# 获取项目根目录的绝对路径
PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "${PROJECT_ROOT}"

# 导入配置和工具函数
source "${PROJECT_ROOT}/scripts/tests/config.sh"

# 定义测试函数
source "${PROJECT_ROOT}/scripts/tests/health_test.sh"
source "${PROJECT_ROOT}/scripts/tests/auth_test.sh"
source "${PROJECT_ROOT}/scripts/tests/user_test.sh"
source "${PROJECT_ROOT}/scripts/tests/social_test.sh"
source "${PROJECT_ROOT}/scripts/tests/profile_test.sh"
source "${PROJECT_ROOT}/scripts/tests/organization_test.sh"

# 检查必要的命令是否存在
if ! command -v jq &> /dev/null; then
    echo "错误: 请先安装 jq"
    exit 1
fi

if ! command -v curl &> /dev/null; then
    echo "错误: 请先安装 curl"
    exit 1
fi

# 主函数
main() {
    log_info "开始 API 测试..."
    
    # 运行健康检查测试
    run_health_tests || {
        log_error "健康检查失败，终止测试"
        exit 1
    }

    # 运行注册测试
    run_register_test

    # 运行认证测试（登录）
    run_auth_tests
    
    # 获取 token 后运行其他测试
    if [ -n "$TOKEN" ]; then
        run_user_tests
        run_profile_tests
        run_profile_metadata_tests
        run_organization_tests
        # 最后运行登出测试
        run_logout_test
    else
        log_error "认证失败，跳过其他测试"
    fi
    
    # 输出测试结果
    print_test_results
    
    # 返回测试结果
    if [ $PASSED_TESTS -eq $TOTAL_TESTS ]; then
        exit 0
    else
        exit 1
    fi
}

# 运行测试
main