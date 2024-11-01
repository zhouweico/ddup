package db

import (
	"testing"
	"time"

	"ddup-apis/internal/config"

	"github.com/stretchr/testify/mock"
)

// 添加测试数据库配置
func setupTestDB() *config.Config {
	return &config.Config{
		Database: struct {
			Host         string
			Port         string
			Name         string
			User         string
			Password     string
			MaxOpenConns int
			MaxIdleConns int
			MaxLifetime  time.Duration
			MaxIdleTime  time.Duration
			RetryTimes   int
			RetryDelay   time.Duration
		}{
			Host:         "localhost",
			Port:         "5432",
			Name:         "test_db",
			User:         "postgres", // 使用默认用户
			Password:     "postgres", // 使用默认密码
			MaxOpenConns: 10,
			MaxIdleConns: 5,
			MaxLifetime:  time.Hour,
			MaxIdleTime:  time.Minute * 5,
			RetryTimes:   3,
			RetryDelay:   time.Second,
		},
	}
}

// 添加 Mock 数据库
type mockDB struct {
	mock.Mock
}

func (m *mockDB) Ping() error {
	args := m.Called()
	return args.Error(0)
}

func TestInitDB(t *testing.T) {
	tests := []struct {
		name    string
		cfg     *config.Config
		wantErr bool
	}{
		{
			name: "正常连接",
			cfg: &config.Config{
				Database: struct {
					Host         string
					Port         string
					Name         string
					User         string
					Password     string
					MaxOpenConns int
					MaxIdleConns int
					MaxLifetime  time.Duration
					MaxIdleTime  time.Duration
					RetryTimes   int
					RetryDelay   time.Duration
				}{
					Host:         "localhost",
					Port:         "5432",
					Name:         "testdb",
					User:         "testuser",
					Password:     "testpass",
					MaxOpenConns: 10,
					MaxIdleConns: 5,
					MaxLifetime:  time.Hour,
					MaxIdleTime:  time.Minute * 30,
					RetryTimes:   3,
					RetryDelay:   time.Second * 2,
				},
			},
			wantErr: false,
		},
		{
			name: "连接失败-错误的主机",
			cfg: &config.Config{
				Database: struct {
					Host         string
					Port         string
					Name         string
					User         string
					Password     string
					MaxOpenConns int
					MaxIdleConns int
					MaxLifetime  time.Duration
					MaxIdleTime  time.Duration
					RetryTimes   int
					RetryDelay   time.Duration
				}{
					Host:       "nonexistent",
					Port:       "5432",
					Name:       "testdb",
					User:       "testuser",
					Password:   "testpass",
					RetryTimes: 1,
					RetryDelay: time.Second,
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := InitDB(tt.cfg)
			if (err != nil) != tt.wantErr {
				t.Errorf("InitDB() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if DB == nil {
					t.Error("InitDB() DB is nil")
				}

				// 测试数据库连接
				sqlDB, err := DB.DB()
				if err != nil {
					t.Errorf("获取数据库实例失败: %v", err)
				}

				// 测试连接池参数
				if max := sqlDB.Stats().MaxOpenConnections; max != tt.cfg.Database.MaxOpenConns {
					t.Errorf("MaxOpenConns = %v, want %v", max, tt.cfg.Database.MaxOpenConns)
				}

				// 测试 Ping
				if err := Ping(); err != nil {
					t.Errorf("Ping() failed: %v", err)
				}

				// 清理
				sqlDB.Close()
			}
		})
	}
}

func TestPing(t *testing.T) {
	tests := []struct {
		name    string
		setup   func()
		wantErr bool
	}{
		{
			name: "数据库已连接",
			setup: func() {
				cfg := &config.Config{
					Database: struct {
						Host         string
						Port         string
						Name         string
						User         string
						Password     string
						MaxOpenConns int
						MaxIdleConns int
						MaxLifetime  time.Duration
						MaxIdleTime  time.Duration
						RetryTimes   int
						RetryDelay   time.Duration
					}{
						Host:     "localhost",
						Port:     "5432",
						Name:     "testdb",
						User:     "testuser",
						Password: "testpass",
					},
				}
				InitDB(cfg)
			},
			wantErr: false,
		},
		{
			name: "数据库未连接",
			setup: func() {
				// 清理数据库连接
				DB = nil
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			err := Ping()
			if (err != nil) != tt.wantErr {
				t.Errorf("Ping() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
