package utils

import (
	"testing"
	"time"

	"ddup-apis/internal/config"
	"ddup-apis/internal/model"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	// 使用测试数据库配置
	dsn := "host=localhost user=testuser password=testpass dbname=testdb port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("设置测试数据库失败: %v", err)
	}

	// 自动迁移测试表
	err = db.AutoMigrate(&model.UserSession{})
	if err != nil {
		t.Fatalf("迁移测试表失败: %v", err)
	}

	return db
}

func TestGenerateToken(t *testing.T) {
	// 设置测试配置
	cfg := &config.Config{}
	cfg.JWT.Secret = "test-secret"
	cfg.JWT.ExpiresIn = time.Hour * 24
	config.SetConfig(*cfg)

	// 设置测试数据库
	db := setupTestDB(t)
	defer func() {
		sqlDB, err := db.DB()
		if err == nil {
			sqlDB.Close()
		}
	}()

	tests := []struct {
		name      string
		userID    uint
		username  string
		setupFunc func()
		wantErr   bool
	}{
		{
			name:     "生成新token",
			userID:   1,
			username: "testuser",
			setupFunc: func() {
				db.Exec("DELETE FROM user_sessions WHERE user_id = ?", 1)
			},
			wantErr: false,
		},
		{
			name:     "已有token时生成新token",
			userID:   1,
			username: "testuser",
			setupFunc: func() {
				session := model.UserSession{
					UserID:    1,
					Token:     "old-token",
					IsValid:   true,
					ExpiredAt: time.Now().Add(time.Hour),
				}
				db.Create(&session)
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setupFunc != nil {
				tt.setupFunc()
			}

			token, createdAt, expiresIn, expiredAt, err := GenerateToken(tt.userID, tt.username)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if token == "" {
					t.Error("GenerateToken() returned empty token")
				}
				if createdAt.IsZero() {
					t.Error("GenerateToken() returned zero createdAt")
				}
				if expiresIn != int64(cfg.JWT.ExpiresIn.Seconds()) {
					t.Errorf("GenerateToken() expiresIn = %v, want %v", expiresIn, cfg.JWT.ExpiresIn.Seconds())
				}
				if expiredAt.IsZero() {
					t.Error("GenerateToken() returned zero expiredAt")
				}
			}
		})
	}
}
