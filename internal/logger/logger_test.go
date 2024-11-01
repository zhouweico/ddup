package logger

import (
	"os"
	"testing"

	"ddup-apis/internal/config"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func TestInitLogger(t *testing.T) {
	tests := []struct {
		name    string
		cfg     *config.Config
		wantErr bool
	}{
		{
			name: "正常初始化",
			cfg: &config.Config{
				Log: struct {
					Level      zapcore.Level `mapstructure:"level"`
					Filename   string        `mapstructure:"filename"`
					MaxSize    int           `mapstructure:"max_size"`
					MaxBackups int           `mapstructure:"max_backups"`
					MaxAge     int           `mapstructure:"max_age"`
					Compress   bool          `mapstructure:"compress"`
				}{
					Level:      zapcore.InfoLevel,
					Filename:   "test.log",
					MaxSize:    10,
					MaxBackups: 5,
					MaxAge:     30,
					Compress:   true,
				},
			},
			wantErr: false,
		},
		{
			name: "无效的日志级别",
			cfg: &config.Config{
				Log: struct {
					Level      zapcore.Level `mapstructure:"level"`
					Filename   string        `mapstructure:"filename"`
					MaxSize    int           `mapstructure:"max_size"`
					MaxBackups int           `mapstructure:"max_backups"`
					MaxAge     int           `mapstructure:"max_age"`
					Compress   bool          `mapstructure:"compress"`
				}{
					Level:    zapcore.Level(99),
					Filename: "test.log",
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 清理之前的日志文件
			os.Remove(tt.cfg.Log.Filename)

			err := InitLogger(tt.cfg)
			if (err != nil) != tt.wantErr {
				t.Errorf("InitLogger() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				// 测试日志记录
				Info("test info message")
				Error("test error message", zap.Error(os.ErrNotExist))
				Debug("test debug message")

				// 验证日志文件是否创建
				if _, err := os.Stat(tt.cfg.Log.Filename); os.IsNotExist(err) {
					t.Error("日志文件未创建")
				}

				// 清理
				os.Remove(tt.cfg.Log.Filename)
			}
		})
	}
}

func TestLoggerMethods(t *testing.T) {
	// 设置测试配置
	cfg := &config.Config{
		Log: struct {
			Level      zapcore.Level `mapstructure:"level"`
			Filename   string        `mapstructure:"filename"`
			MaxSize    int           `mapstructure:"max_size"`
			MaxBackups int           `mapstructure:"max_backups"`
			MaxAge     int           `mapstructure:"max_age"`
			Compress   bool          `mapstructure:"compress"`
		}{
			Level:    zapcore.DebugLevel,
			Filename: "test_methods.log",
		},
	}

	if err := InitLogger(cfg); err != nil {
		t.Fatalf("初始化日志失败: %v", err)
	}
	defer os.Remove(cfg.Log.Filename)
}
