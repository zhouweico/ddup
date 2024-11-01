package config

import (
	"os"
	"testing"
	"time"

	"go.uber.org/zap/zapcore"
)

func TestLoadConfig(t *testing.T) {
	// 创建临时配置文件
	tmpfile, err := os.CreateTemp("", "test.*.env")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	// 写入测试配置
	content := `
SERVER_PORT=8080
SERVER_MODE=debug
DATABASE_HOST=localhost
DATABASE_PORT=5432
DATABASE_NAME=testdb
DATABASE_USER=testuser
DATABASE_PASSWORD=testpass
DATABASE_MAX_OPEN_CONNS=10
DATABASE_MAX_IDLE_CONNS=5
DATABASE_MAX_LIFETIME=1h
DATABASE_MAX_IDLE_TIME=30m
DATABASE_RETRY_TIMES=3
DATABASE_RETRY_DELAY=5s
JWT_SECRET=test-secret
JWT_EXPIRES_IN=24h
HEALTH_CHECK_INTERVAL=30s
LOG_LEVEL=info
LOG_FILENAME=test.log
LOG_MAX_SIZE=10
LOG_MAX_BACKUPS=5
LOG_MAX_AGE=30
LOG_COMPRESS=true
`
	if err := os.WriteFile(tmpfile.Name(), []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	// 设置环境变量
	os.Setenv("ENV_FILE", tmpfile.Name())
	defer os.Unsetenv("ENV_FILE")

	// 测试加载配置
	cfg, err := LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig() error = %v", err)
	}

	// 验证配置值
	tests := []struct {
		name     string
		got      interface{}
		want     interface{}
		errorMsg string
	}{
		{"Server.Port", cfg.Server.Port, "8080", "服务器端口不匹配"},
		{"Server.Mode", cfg.Server.Mode, "debug", "服务器模式不匹配"},
		{"Database.Host", cfg.Database.Host, "localhost", "数据库主机不匹配"},
		{"Database.MaxOpenConns", cfg.Database.MaxOpenConns, 10, "最大连接数不匹配"},
		{"JWT.Secret", cfg.JWT.Secret, "test-secret", "JWT密钥不匹配"},
		{"JWT.ExpiresIn", cfg.JWT.ExpiresIn, 24 * time.Hour, "JWT过期时间不匹配"},
		{"Log.Level", cfg.Log.Level, zapcore.InfoLevel, "日志级别不匹配"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.want {
				t.Errorf("%s = %v, want %v", tt.errorMsg, tt.got, tt.want)
			}
		})
	}
}

func TestGetConfig(t *testing.T) {
	t.Run("未初始化配置", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("期望在配置未初始化时触发 panic")
			}
		}()
		GetConfig()
	})

	t.Run("已初始化配置", func(t *testing.T) {
		// 设置测试配置
		globalConfig = Config{
			JWT: struct {
				Secret    string
				ExpiresIn time.Duration
			}{
				Secret: "test-secret",
			},
		}

		cfg := GetConfig()
		if cfg.JWT.Secret != "test-secret" {
			t.Errorf("GetConfig().JWT.Secret = %v, want %v", cfg.JWT.Secret, "test-secret")
		}
	})
}
