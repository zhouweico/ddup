# DDUP

## 项目说明
基于 Golang 的 RESTful API 服务，使用 PostgreSQL 作为数据存储。

## 已实现功能

### 用户管理
- [x] 用户注册
- [x] 用户登录
- [x] 用户退出
- [x] JWT 认证
- [x] 获取用户详情
- [x] 更新用户信息
- [x] 修改密码
- [x] 删除用户（软删除）

### 安全特性
- [x] JWT Token 认证
- [x] 资源访问控制
- [x] 密码加密存储
- [x] CORS 跨域支持

### 系统特性
- [x] RESTful API 设计
- [x] Swagger API 文档
- [x] 统一响应格式
- [x] 分页查询支持
- [x] 环境配置管理
- [x] 数据库连接池
- [x] 数据库重试机制
- [x] 健康检查
- [x] 日志系统
  - 同时输出到控制台和文件
  - 日志轮转
  - 支持不同日志级别
  - 结构化日志输出
  - 自动记录调用位置

## 技术栈
- Go 1.21+
- PostgreSQL 14+
- Gin Web Framework
- GORM

## 开发环境搭建
1. 数据库环境
```bash
docker run -d \
    --name ddup-postgres \
    -p 5432:5432 \
    -e POSTGRES_PASSWORD=Admin@123456 \
    -e PGDATA=/var/lib/postgresql/data/pgdata \
    -v $HOME/data/ddup/data:/var/lib/postgresql/data \
    postgres:17.0
```

2. 创建数据库和用户

```bash
# 1. 创建用户并设置密码
docker exec -it ddup-postgres psql -U postgres -c "CREATE USER ddup WITH PASSWORD 'Ddup@123456';"

# 2. 创建数据库
docker exec -it ddup-postgres psql -U postgres -c "CREATE DATABASE ddup owner ddup;"


# 3. 授予权限
docker exec -it ddup-postgres psql -U postgres -c "GRANT ALL PRIVILEGES ON DATABASE ddup TO ddup;"
docker exec -it ddup-postgres psql -U postgres -d ddup -c "GRANT ALL ON SCHEMA public TO ddup;"
```

3. 复制环境变量文件

```bash
cp .env.example .env
```

4. 修改环境变量配置

5. 运行服务
```bash
go run cmd/api/main.go
```

## 项目结构
```
.
├── cmd/                # 主要的应用程序入口
│   └── api/           # API 服务入口
├── configs/           # 配置文件
├── internal/          # 私有应用程序和库代码
│   ├── config/       # 配置
│   ├── handler/      # HTTP 处理器
│   ├── middleware/   # HTTP 中间件
│   ├── model/        # 数据库模型
│   ├── repository/   # 数据库操作
│   └── service/      # 业务逻辑
├── pkg/              # 可以被外部应用程序使用的库代码
│   ├── logger/       # 日志工具
│   └── utils/        # 通用工具
└── scripts/          # 脚本和工具
```

## 开发工具

### 提交代码

使用提供的脚本自动生成 Swagger 文档并提交代码：
```bash
./scripts/commit.sh "提交信息"
```

## 配置说明

主要配置项（.env 文件）：

- 服务配置：端口、环境等
- 数据库配置：连接信息、连接池参数等
- JWT 配置：密钥、过期时间等
- 健康检查配置：检查间隔等
- 日志配置：
  - 日志级别
  - 日志文件路径
  - 日志文件大小限制
  - 日志文件保留数量
  - 日志文件保留天数
  - 是否压缩

## 贡献

欢迎提交 Issue 和 Pull Request

## 许可证

MIT License