# 构建阶段
FROM golang:1.23.2-alpine AS builder

# 设置工作目录
WORKDIR /app

# 安装依赖
RUN apk add --no-cache git

# 复制 go.mod 和 go.sum
COPY go.mod go.sum ./
RUN go mod download

# 复制源代码
COPY . .

# 构建应用
RUN CGO_ENABLED=0 GOOS=linux go build -o ddup-api cmd/api/main.go

# 运行阶段
FROM alpine:latest

# 安装基础工具
RUN apk --no-cache add ca-certificates tzdata

# 设置时区
ENV TZ=Asia/Shanghai

WORKDIR /app

# 从构建阶段复制二进制文件
COPY --from=builder /app/ddup-api .
COPY --from=builder /app/.env.example .env

# 暴露端口
EXPOSE 8080

# 启动应用
CMD ["./ddup-api"] 