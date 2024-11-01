package middleware

import (
	"context"
	"testing"
	"time"

	"ddup-apis/internal/config"
	"ddup-apis/internal/db"
)

func TestPeriodicHealthCheck(t *testing.T) {
	// 创建测试通道用于控制测试时间
	done := make(chan bool)

	tests := []struct {
		name     string
		interval time.Duration
		setupDB  func()
		wantErr  bool
	}{
		{
			name:     "正常健康检查",
			interval: 100 * time.Millisecond,
			setupDB: func() {
				// 模拟正常的数据库连接
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
				db.InitDB(cfg)
			},
			wantErr: false,
		},
		{
			name:     "数据库连接失败",
			interval: 100 * time.Millisecond,
			setupDB: func() {
				db.DB = nil
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setupDB != nil {
				tt.setupDB()
			}

			// 创建上下文用于控制测试时间
			ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
			defer cancel()

			// 启动健康检查
			go PeriodicHealthCheck(tt.interval)

			// 等待几个检查周期
			select {
			case <-ctx.Done():
				// 测试完成
			case <-done:
				// 提前结束测试
			}

			// 验证健康检查状态
			if tt.wantErr {
				if IsHealthy() {
					t.Error("期望健康检查失败，但状态为健康")
				}
			} else {
				if !IsHealthy() {
					t.Error("期望健康检查成功，但状态为不健康")
				}
			}
		})
	}
}

func TestIsHealthy(t *testing.T) {
	tests := []struct {
		name    string
		setup   func()
		want    bool
		cleanup func()
	}{
		{
			name: "初始状态",
			setup: func() {
				healthStatus = true
			},
			want: true,
		},
		{
			name: "不健康状态",
			setup: func() {
				healthStatus = false
			},
			want: false,
		},
		{
			name: "状态切换",
			setup: func() {
				healthStatus = true
				go func() {
					time.Sleep(50 * time.Millisecond)
					healthStatus = false
				}()
			},
			want: false,
			cleanup: func() {
				time.Sleep(100 * time.Millisecond)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup()
			}

			if tt.cleanup != nil {
				defer tt.cleanup()
			}

			if got := IsHealthy(); got != tt.want {
				t.Errorf("IsHealthy() = %v, want %v", got, tt.want)
			}
		})
	}
}
