package service

import (
	"context"
	"testing"
	"time"

	"ddup-apis/internal/model"
	"ddup-apis/internal/utils"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type mockDB struct {
	*gorm.DB
	findErr   error
	createErr error
	updateErr error
	deleteErr error
}

func setupTestService(t *testing.T) (*UserService, *mockDB) {
	// 创建测试数据库连接
	dsn := "host=localhost user=testuser password=testpass dbname=testdb port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("设置测试数据库失败: %v", err)
	}

	mock := &mockDB{DB: db}
	service := NewUserService(db)

	return service, mock
}

func TestUserService_Register(t *testing.T) {
	service, mock := setupTestService(t)

	tests := []struct {
		name      string
		user      *model.User
		setupMock func()
		wantErr   bool
	}{
		{
			name: "注册成功",
			user: &model.User{
				Username: "testuser",
				Password: "testpass123",
				Email:    "test@example.com",
			},
			setupMock: func() {
				mock.findErr = gorm.ErrRecordNotFound
				mock.createErr = nil
			},
			wantErr: false,
		},
		{
			name: "用户名已存在",
			user: &model.User{
				Username: "existinguser",
				Password: "testpass123",
				Email:    "test@example.com",
			},
			setupMock: func() {
				mock.findErr = nil
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setupMock != nil {
				tt.setupMock()
			}

			err := service.Register(context.Background(), tt.user.Username, tt.user.Password)
			if (err != nil) != tt.wantErr {
				t.Errorf("UserService.Register() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUserService_Login(t *testing.T) {
	service, mock := setupTestService(t)

	// 创建测试用户
	hashedPassword, _ := utils.HashPassword("testpass123")
	testUser := &model.User{
		Username: "testuser",
		Password: hashedPassword,
	}

	tests := []struct {
		name      string
		username  string
		password  string
		setupMock func()
		wantErr   bool
	}{
		{
			name:     "登录成功",
			username: "testuser",
			password: "testpass123",
			setupMock: func() {
				mock.findErr = nil
				mock.DB.Create(testUser)
			},
			wantErr: false,
		},
		{
			name:     "用户不存在",
			username: "nonexistent",
			password: "testpass123",
			setupMock: func() {
				mock.findErr = gorm.ErrRecordNotFound
			},
			wantErr: true,
		},
		{
			name:     "密码错误",
			username: "testuser",
			password: "wrongpass",
			setupMock: func() {
				mock.findErr = nil
				mock.DB.Create(testUser)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setupMock != nil {
				tt.setupMock()
			}

			loginResult, err := service.Login(context.Background(), tt.username, tt.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("UserService.Login() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if loginResult.User == nil {
					t.Error("UserService.Login() returned nil user")
				}
				if loginResult.Token == "" {
					t.Error("UserService.Login() returned empty token")
				}
			}
		})
	}
}

func TestUserService_GetUserByID(t *testing.T) {
	service, mock := setupTestService(t)

	testUser := &model.User{
		Model:    gorm.Model{ID: 1},
		Username: "testuser",
		Email:    "test@example.com",
	}

	tests := []struct {
		name      string
		userID    uint
		setupMock func()
		want      *model.User
		wantErr   bool
	}{
		{
			name:   "获取成功",
			userID: 1,
			setupMock: func() {
				mock.findErr = nil
				mock.DB.Create(testUser)
			},
			want:    testUser,
			wantErr: false,
		},
		{
			name:   "用户不存在",
			userID: 999,
			setupMock: func() {
				mock.findErr = gorm.ErrRecordNotFound
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setupMock != nil {
				tt.setupMock()
			}

			got, err := service.GetUserByID(context.Background(), tt.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("UserService.GetUserByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && got.ID != tt.want.ID {
				t.Errorf("UserService.GetUserByID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserService_ValidateToken(t *testing.T) {
	service, mock := setupTestService(t)

	// 创建测试用户和会话
	testUser := &model.User{
		Model:    gorm.Model{ID: 1},
		Username: "testuser",
	}
	validSession := &model.UserSession{
		Token:     "valid_token",
		UserID:    testUser.ID,
		IsValid:   true,
		ExpiredAt: time.Now().Add(24 * time.Hour), // 24小时后过期
	}
	expiredSession := &model.UserSession{
		Token:     "expired_token",
		UserID:    testUser.ID,
		IsValid:   true,
		ExpiredAt: time.Now().Add(-24 * time.Hour), // 24小时前过期
	}

	tests := []struct {
		name      string
		token     string
		setupMock func()
		want      *TokenValidationResult
		wantErr   bool
	}{
		{
			name:  "有效token",
			token: "valid_token",
			setupMock: func() {
				mock.DB.Create(testUser)
				mock.DB.Create(validSession)
			},
			want: &TokenValidationResult{
				UserID:   testUser.ID,
				Username: testUser.Username,
				IsValid:  true,
			},
			wantErr: false,
		},
		{
			name:  "过期token",
			token: "expired_token",
			setupMock: func() {
				mock.DB.Create(testUser)
				mock.DB.Create(expiredSession)
			},
			want: &TokenValidationResult{
				IsValid: false,
			},
			wantErr: false,
		},
		{
			name:  "不存在的token",
			token: "nonexistent_token",
			setupMock: func() {
				// 不需要设置任何数据
			},
			want: &TokenValidationResult{
				IsValid: false,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setupMock != nil {
				tt.setupMock()
			}

			got, err := service.ValidateToken(tt.token)
			if (err != nil) != tt.wantErr {
				t.Errorf("UserService.ValidateToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if got.IsValid != tt.want.IsValid {
					t.Errorf("UserService.ValidateToken() IsValid = %v, want %v", got.IsValid, tt.want.IsValid)
				}
				if got.IsValid {
					if got.UserID != tt.want.UserID {
						t.Errorf("UserService.ValidateToken() UserID = %v, want %v", got.UserID, tt.want.UserID)
					}
					if got.Username != tt.want.Username {
						t.Errorf("UserService.ValidateToken() Username = %v, want %v", got.Username, tt.want.Username)
					}
				}
			}
		})
	}
}
