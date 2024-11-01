package service

import (
	"context"
	"reflect"
	"sync"
	"testing"
	"time"

	"ddup-apis/internal/errors"
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

func (m *mockDB) First(dest interface{}, conds ...interface{}) *gorm.DB {
	if m.findErr != nil {
		return &gorm.DB{Error: m.findErr}
	}
	return m.DB.First(dest, conds...)
}

func (m *mockDB) Create(value interface{}) *gorm.DB {
	if m.createErr != nil {
		return &gorm.DB{Error: m.createErr}
	}
	return m.DB.Create(value)
}

func (m *mockDB) Delete(value interface{}, conds ...interface{}) *gorm.DB {
	if m.deleteErr != nil {
		return &gorm.DB{Error: m.deleteErr}
	}
	return m.DB.Delete(value, conds...)
}

func (m *mockDB) Model(value interface{}) *gorm.DB {
	return m.DB.Model(value)
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
	service, _ := setupTestService(t)

	tests := []struct {
		name    string
		token   string
		setup   func()
		want    *TokenValidationResult
		wantErr bool
	}{
		{
			name:  "有效token",
			token: "valid-token",
			setup: func() {
				// 在数据库中创建有效的session
				session := model.UserSession{
					UserID:    1,
					Token:     "valid-token",
					IsValid:   true,
					ExpiredAt: time.Now().Add(time.Hour),
				}
				service.db.Create(&session)
			},
			want: &TokenValidationResult{
				UserID:  1,
				IsValid: true,
			},
			wantErr: false,
		},
		{
			name:    "无效token",
			token:   "invalid-token",
			setup:   func() {},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup()
			}

			got, err := service.ValidateToken(tt.token)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ValidateToken() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserService_UpdateUser(t *testing.T) {
	service, mock := setupTestService(t)

	tests := []struct {
		name    string
		userID  uint
		updates map[string]interface{}
		setup   func()
		wantErr bool
	}{
		{
			name:   "更新成功",
			userID: 1,
			updates: map[string]interface{}{
				"nickname": "新昵称",
				"email":    "new@example.com",
			},
			setup: func() {
				mock.DB.Create(&model.User{
					Model:    gorm.Model{ID: 1},
					Username: "testuser",
				})
			},
			wantErr: false,
		},
		{
			name:   "用户不存在",
			userID: 999,
			updates: map[string]interface{}{
				"nickname": "新昵称",
			},
			setup:   func() {},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup()
			}

			err := service.UpdateUser(context.Background(), tt.userID, tt.updates)
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateUser() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUserService_DeleteUser(t *testing.T) {
	service, mock := setupTestService(t)

	tests := []struct {
		name    string
		userID  uint
		setup   func()
		wantErr bool
	}{
		{
			name:   "删除成功",
			userID: 1,
			setup: func() {
				user := &model.User{
					Model:    gorm.Model{ID: 1},
					Username: "testuser",
					Email:    "test@example.com",
				}
				mock.DB.Create(user)
			},
			wantErr: false,
		},
		{
			name:   "用户不存在",
			userID: 999,
			setup: func() {
				mock.findErr = gorm.ErrRecordNotFound
			},
			wantErr: true,
		},
		{
			name:   "数据库错误",
			userID: 1,
			setup: func() {
				mock.deleteErr = errors.New(500, "数据库错误", nil)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup()
			}

			err := service.DeleteUser(context.Background(), tt.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteUser() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr {
				// 验证用户是否已被软删除
				var user model.User
				err := mock.DB.Unscoped().First(&user, tt.userID).Error
				if err != nil {
					t.Errorf("查找删除的用户失败: %v", err)
				}
				if user.DeletedAt.Time.IsZero() {
					t.Error("用户未被软删除")
				}
			}
		})
	}
}

func TestUserService_Logout(t *testing.T) {
	service, mock := setupTestService(t)

	tests := []struct {
		name    string
		token   string
		setup   func()
		wantErr bool
	}{
		{
			name:  "登出成功",
			token: "valid-token",
			setup: func() {
				session := &model.UserSession{
					UserID:    1,
					Token:     "valid-token",
					IsValid:   true,
					ExpiredAt: time.Now().Add(time.Hour),
				}
				mock.DB.Create(session)
			},
			wantErr: false,
		},
		{
			name:  "无效token",
			token: "invalid-token",
			setup: func() {
				mock.findErr = gorm.ErrRecordNotFound
			},
			wantErr: true,
		},
		{
			name:  "已登出的token",
			token: "logged-out-token",
			setup: func() {
				session := &model.UserSession{
					UserID:    1,
					Token:     "logged-out-token",
					IsValid:   false,
					ExpiredAt: time.Now().Add(time.Hour),
				}
				mock.DB.Create(session)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup()
			}

			err := service.Logout(context.Background(), tt.token)
			if (err != nil) != tt.wantErr {
				t.Errorf("Logout() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr {
				// 验证session是否已被标记为无效
				var session model.UserSession
				err := mock.DB.Where("token = ?", tt.token).First(&session).Error
				if err != nil {
					t.Errorf("查找session失败: %v", err)
				}
				if session.IsValid {
					t.Error("session未被标记为无效")
				}
			}
		})
	}
}

// 添加一些边界情况的测试
func TestUserService_EdgeCases(t *testing.T) {
	service, mock := setupTestService(t)

	t.Run("并发登录测试", func(t *testing.T) {
		userID := uint(1)
		username := "testuser"

		// 创建测试用户
		mock.DB.Create(&model.User{
			Model:    gorm.Model{ID: userID},
			Username: username,
			Password: "hashedpass",
		})

		// 模拟并发登录
		var wg sync.WaitGroup
		for i := 0; i < 5; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				_, err := service.Login(context.Background(), username, "password")
				if err != nil {
					t.Errorf("并发登录失败: %v", err)
				}
			}()
		}
		wg.Wait()

		// 验证只有最新的session是有效的
		var sessions []model.UserSession
		mock.DB.Where("user_id = ?", userID).Find(&sessions)

		validCount := 0
		for _, session := range sessions {
			if session.IsValid {
				validCount++
			}
		}

		if validCount != 1 {
			t.Errorf("有效session数量 = %d, 期望 1", validCount)
		}
	})

	t.Run("密码重试限制测试", func(t *testing.T) {
		user := &model.User{
			Model:         gorm.Model{ID: 2},
			Username:      "testuser2",
			Password:      "hashedpass",
			LoginAttempts: 0,
		}
		mock.DB.Create(user)

		// 尝试多次错误登录
		for i := 0; i < 5; i++ {
			service.Login(context.Background(), user.Username, "wrongpass")
		}

		// 验证账户是否被锁定
		var updatedUser model.User
		mock.DB.First(&updatedUser, user.ID)
		if updatedUser.LoginAttempts < 5 {
			t.Error("登录尝试次数未正确记录")
		}
		if updatedUser.LockedUntil == nil || updatedUser.LockedUntil.Before(time.Now()) {
			t.Error("账户未被锁定")
		}
	})

	t.Run("Token过期测试", func(t *testing.T) {
		expiredToken := "expired-token"
		session := &model.UserSession{
			UserID:    3,
			Token:     expiredToken,
			IsValid:   true,
			ExpiredAt: time.Now().Add(-time.Hour), // 已过期
		}
		mock.DB.Create(session)

		result, err := service.ValidateToken(expiredToken)
		if err == nil {
			t.Error("期望过期token返回错误")
		}
		if result != nil {
			t.Error("过期token不应返回有效结果")
		}
	})
}
