# DDUP

## 项目说明
基于 Golang 的 RESTful API 服务，使用 PostgreSQL 作为数据存储。

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
# 1. 创建数据库
docker exec -it ddup-postgres psql -U postgres -c "CREATE DATABASE ddup;"

# 2. 创建用户并设置密码
docker exec -it ddup-postgres psql -U postgres -c "CREATE USER ddup WITH PASSWORD 'Ddup@123456';"

# 3. 授予权限
docker exec -it ddup-postgres psql -U postgres -c "GRANT ALL PRIVILEGES ON DATABASE ddup TO ddup;"
docker exec -it ddup-postgres psql -U postgres -d ddup -c "GRANT ALL ON SCHEMA public TO ddup;"

# 4. 设置数据库所有者
docker exec -it ddup-postgres psql -U postgres -c "ALTER DATABASE ddup OWNER TO ddup;"
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
