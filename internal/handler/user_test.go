package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"ddup-apis/internal/model"
	"ddup-apis/internal/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var _ service.IUserService = (*mockUserService)(nil)

type mockUserService struct {
	registerErr    error
	loginResult    *service.LoginResult
	loginErr       error
	getUser        *model.User
	getUserErr     error
	updateErr      error
	deleteErr      error
	changePassErr  error
	logoutErr      error
	user           *model.User
	validateResult *service.TokenValidationResult
}

func (m *mockUserService) Register(ctx context.Context, username, password string) error {
	return m.registerErr
}

func (m *mockUserService) Login(ctx context.Context, username, password string) (*service.LoginResult, error) {
	if m.loginErr != nil {
		return nil, m.loginErr
	}
	return m.loginResult, nil
}

func (m *mockUserService) ValidateToken(token string) (*service.TokenValidationResult, error) {
	return nil, nil
}

func (m *mockUserService) GetUserByID(ctx context.Context, userID uint) (*model.User, error) {
	return m.getUser, m.getUserErr
}

func (m *mockUserService) UpdateUser(ctx context.Context, userID uint, updates map[string]interface{}) error {
	return m.updateErr
}

func (m *mockUserService) DeleteUser(ctx context.Context, userID uint) error {
	return m.deleteErr
}

func (m *mockUserService) ChangePassword(ctx context.Context, userID uint, oldPassword, newPassword string) error {
	return m.changePassErr
}

func (m *mockUserService) Logout(ctx context.Context, token string) error {
	return m.logoutErr
}

func setupTestHandler() (*UserHandler, *gin.Engine, *mockUserService) {
	gin.SetMode(gin.TestMode)
	mock := &mockUserService{}
	handler := NewUserHandler(mock)
	router := gin.New()
	return handler, router, mock
}

func TestUserHandler_Register(t *testing.T) {
	handler, router, mock := setupTestHandler()
	router.POST("/register", handler.Register)

	tests := []struct {
		name       string
		reqBody    map[string]interface{}
		setupMock  func()
		wantStatus int
		wantBody   map[string]interface{}
	}{
		{
			name: "注册成功",
			reqBody: map[string]interface{}{
				"username": "testuser",
				"password": "testpass123",
				"email":    "test@example.com",
			},
			setupMock: func() {
				mock.registerErr = nil
			},
			wantStatus: http.StatusOK,
			wantBody: map[string]interface{}{
				"code":    float64(http.StatusOK),
				"message": "注册成功",
			},
		},
		{
			name: "用户名已存在",
			reqBody: map[string]interface{}{
				"username": "existinguser",
				"password": "testpass123",
				"email":    "test@example.com",
			},
			setupMock: func() {
				mock.registerErr = errors.New("用户名已存在")
			},
			wantStatus: http.StatusBadRequest,
			wantBody: map[string]interface{}{
				"code":    float64(http.StatusBadRequest),
				"message": "用户名已存在",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setupMock != nil {
				tt.setupMock()
			}

			body, _ := json.Marshal(tt.reqBody)
			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/register", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			router.ServeHTTP(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("Register() status = %v, want %v", w.Code, tt.wantStatus)
			}

			var got map[string]interface{}
			json.Unmarshal(w.Body.Bytes(), &got)

			if got["code"] != tt.wantBody["code"] {
				t.Errorf("Register() code = %v, want %v", got["code"], tt.wantBody["code"])
			}
			if got["message"] != tt.wantBody["message"] {
				t.Errorf("Register() message = %v, want %v", got["message"], tt.wantBody["message"])
			}
		})
	}
}

func TestUserHandler_Login(t *testing.T) {
	handler, router, mock := setupTestHandler()
	router.POST("/login", handler.Login)

	tests := []struct {
		name       string
		reqBody    map[string]interface{}
		setupMock  func()
		wantStatus int
		wantBody   map[string]interface{}
	}{
		{
			name: "登录成功",
			reqBody: map[string]interface{}{
				"username": "testuser",
				"password": "testpass123",
			},
			setupMock: func() {
				mock.loginResult = &service.LoginResult{
					Token:     "test-token",
					User:      &model.User{Username: "testuser"},
					ExpiresIn: 3600,
				}
				mock.loginErr = nil
			},
			wantStatus: http.StatusOK,
			wantBody: map[string]interface{}{
				"code":    float64(http.StatusOK),
				"message": "登录成功",
				"data": map[string]interface{}{
					"token": "test-token",
				},
			},
		},
		{
			name: "用户名或密码错误",
			reqBody: map[string]interface{}{
				"username": "wronguser",
				"password": "wrongpass",
			},
			setupMock: func() {
				mock.loginErr = errors.New("用户名或密码错误")
			},
			wantStatus: http.StatusUnauthorized,
			wantBody: map[string]interface{}{
				"code":    float64(http.StatusUnauthorized),
				"message": "用户名���密码错误",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setupMock != nil {
				tt.setupMock()
			}

			body, _ := json.Marshal(tt.reqBody)
			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/login", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			router.ServeHTTP(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("Login() status = %v, want %v", w.Code, tt.wantStatus)
			}

			var got map[string]interface{}
			json.Unmarshal(w.Body.Bytes(), &got)

			if got["code"] != tt.wantBody["code"] {
				t.Errorf("Login() code = %v, want %v", got["code"], tt.wantBody["code"])
			}
			if got["message"] != tt.wantBody["message"] {
				t.Errorf("Login() message = %v, want %v", got["message"], tt.wantBody["message"])
			}
		})
	}
}

func TestUserHandler_GetUser(t *testing.T) {
	handler, router, mock := setupTestHandler()
	router.GET("/user", handler.GetUser)

	tests := []struct {
		name       string
		setupMock  func()
		setupAuth  func(*http.Request)
		wantStatus int
		wantBody   map[string]interface{}
	}{
		{
			name: "获取用户成功",
			setupMock: func() {
				mock.getUser = &model.User{
					Username: "testuser",
					Email:    "test@example.com",
				}
				mock.getUserErr = nil
			},
			setupAuth: func(req *http.Request) {
				req.Header.Set("Authorization", "Bearer valid-token")
				// 模拟中间件设置用户ID
				ctx := req.Context()
				ctx = context.WithValue(ctx, "userID", uint(1))
				*req = *req.WithContext(ctx)
			},
			wantStatus: http.StatusOK,
			wantBody: map[string]interface{}{
				"code":    float64(http.StatusOK),
				"message": "获取用户信息成功",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setupMock != nil {
				tt.setupMock()
			}

			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/user", nil)
			if tt.setupAuth != nil {
				tt.setupAuth(req)
			}

			router.ServeHTTP(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("GetUser() status = %v, want %v", w.Code, tt.wantStatus)
			}

			var got map[string]interface{}
			json.Unmarshal(w.Body.Bytes(), &got)

			if got["code"] != tt.wantBody["code"] {
				t.Errorf("GetUser() code = %v, want %v", got["code"], tt.wantBody["code"])
			}
			if got["message"] != tt.wantBody["message"] {
				t.Errorf("GetUser() message = %v, want %v", got["message"], tt.wantBody["message"])
			}
		})
	}
}

func TestUserHandler_UpdateUser(t *testing.T) {
	handler, router, mock := setupTestHandler()
	router.PUT("/user", handler.UpdateUser)

	tests := []struct {
		name       string
		setupAuth  func(*http.Request)
		reqBody    map[string]interface{}
		setupMock  func()
		wantStatus int
		wantBody   map[string]interface{}
	}{
		{
			name: "更新成功",
			setupAuth: func(req *http.Request) {
				req.Header.Set("Authorization", "Bearer valid-token")
				// 模拟认证中间件设置用户ID
				router.Use(func(c *gin.Context) {
					c.Set("userID", uint(1))
					c.Next()
				})
			},
			reqBody: map[string]interface{}{
				"nickname": "新昵称",
				"email":    "new@example.com",
			},
			setupMock: func() {
				mock.getUserErr = nil
				mock.user = &model.User{
					Model: gorm.Model{
						ID: 1,
					},
					Username: "testuser",
				}
				mock.updateErr = nil
			},
			wantStatus: http.StatusOK,
			wantBody: map[string]interface{}{
				"code":    float64(http.StatusOK),
				"message": "更新成功",
			},
		},
		{
			name: "未授权访问",
			setupAuth: func(req *http.Request) {
				// 不设置认证信息
			},
			reqBody: map[string]interface{}{
				"nickname": "新昵称",
			},
			setupMock:  func() {},
			wantStatus: http.StatusUnauthorized,
			wantBody: map[string]interface{}{
				"code":    float64(http.StatusUnauthorized),
				"message": "未授权",
			},
		},
		{
			name: "用户不存在",
			setupAuth: func(req *http.Request) {
				req.Header.Set("Authorization", "Bearer valid-token")
				router.Use(func(c *gin.Context) {
					c.Set("userID", uint(999))
					c.Next()
				})
			},
			reqBody: map[string]interface{}{
				"nickname": "新昵称",
			},
			setupMock: func() {
				mock.getUserErr = errors.New("用户不存在")
			},
			wantStatus: http.StatusNotFound,
			wantBody: map[string]interface{}{
				"code":    float64(http.StatusNotFound),
				"message": "用户不存在",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setupMock != nil {
				tt.setupMock()
			}

			body, _ := json.Marshal(tt.reqBody)
			w := httptest.NewRecorder()
			req := httptest.NewRequest("PUT", "/user", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			if tt.setupAuth != nil {
				tt.setupAuth(req)
			}

			router.ServeHTTP(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("UpdateUser() status = %v, want %v", w.Code, tt.wantStatus)
			}

			var got map[string]interface{}
			json.Unmarshal(w.Body.Bytes(), &got)

			if got["code"] != tt.wantBody["code"] {
				t.Errorf("UpdateUser() code = %v, want %v", got["code"], tt.wantBody["code"])
			}
			if got["message"] != tt.wantBody["message"] {
				t.Errorf("UpdateUser() message = %v, want %v", got["message"], tt.wantBody["message"])
			}
		})
	}
}

func TestUserHandler_DeleteUser(t *testing.T) {
	handler, router, mock := setupTestHandler()
	router.DELETE("/user", handler.DeleteUser)

	tests := []struct {
		name       string
		setupAuth  func(*http.Request)
		setupMock  func()
		wantStatus int
		wantBody   map[string]interface{}
	}{
		{
			name: "删除成功",
			setupAuth: func(req *http.Request) {
				req.Header.Set("Authorization", "Bearer valid-token")
				router.Use(func(c *gin.Context) {
					c.Set("userID", uint(1))
					c.Next()
				})
			},
			setupMock: func() {
				mock.getUserErr = nil
				mock.user = &model.User{
					Model: gorm.Model{
						ID: 1,
					},
					Username: "testuser",
				}
				mock.deleteErr = nil
			},
			wantStatus: http.StatusOK,
			wantBody: map[string]interface{}{
				"code":    float64(http.StatusOK),
				"message": "删除成功",
			},
		},
		{
			name: "未授权访问",
			setupAuth: func(req *http.Request) {
				// 不设置认证信息
			},
			setupMock:  func() {},
			wantStatus: http.StatusUnauthorized,
			wantBody: map[string]interface{}{
				"code":    float64(http.StatusUnauthorized),
				"message": "未授权",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setupMock != nil {
				tt.setupMock()
			}

			w := httptest.NewRecorder()
			req := httptest.NewRequest("DELETE", "/user", nil)

			if tt.setupAuth != nil {
				tt.setupAuth(req)
			}

			router.ServeHTTP(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("DeleteUser() status = %v, want %v", w.Code, tt.wantStatus)
			}

			var got map[string]interface{}
			json.Unmarshal(w.Body.Bytes(), &got)

			if got["code"] != tt.wantBody["code"] {
				t.Errorf("DeleteUser() code = %v, want %v", got["code"], tt.wantBody["code"])
			}
			if got["message"] != tt.wantBody["message"] {
				t.Errorf("DeleteUser() message = %v, want %v", got["message"], tt.wantBody["message"])
			}
		})
	}
}

func TestUserHandler_ChangePassword(t *testing.T) {
	handler, router, mock := setupTestHandler()
	router.PUT("/user/password", handler.ChangePassword)

	tests := []struct {
		name       string
		setupAuth  func(*http.Request)
		reqBody    map[string]interface{}
		setupMock  func()
		wantStatus int
		wantBody   map[string]interface{}
	}{
		{
			name: "修改密码成功",
			setupAuth: func(req *http.Request) {
				req.Header.Set("Authorization", "Bearer valid-token")
				router.Use(func(c *gin.Context) {
					c.Set("userID", uint(1))
					c.Next()
				})
			},
			reqBody: map[string]interface{}{
				"old_password": "oldpass123",
				"new_password": "newpass123",
			},
			setupMock: func() {
				mock.changePassErr = nil
			},
			wantStatus: http.StatusOK,
			wantBody: map[string]interface{}{
				"code":    float64(http.StatusOK),
				"message": "密码修改成功",
			},
		},
		{
			name: "原密码错误",
			setupAuth: func(req *http.Request) {
				req.Header.Set("Authorization", "Bearer valid-token")
				router.Use(func(c *gin.Context) {
					c.Set("userID", uint(1))
					c.Next()
				})
			},
			reqBody: map[string]interface{}{
				"old_password": "wrongpass",
				"new_password": "newpass123",
			},
			setupMock: func() {
				mock.changePassErr = errors.New("原密码错误")
			},
			wantStatus: http.StatusBadRequest,
			wantBody: map[string]interface{}{
				"code":    float64(http.StatusBadRequest),
				"message": "原密码错误",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setupMock != nil {
				tt.setupMock()
			}

			body, _ := json.Marshal(tt.reqBody)
			w := httptest.NewRecorder()
			req := httptest.NewRequest("PUT", "/user/password", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			if tt.setupAuth != nil {
				tt.setupAuth(req)
			}

			router.ServeHTTP(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("ChangePassword() status = %v, want %v", w.Code, tt.wantStatus)
			}

			var got map[string]interface{}
			json.Unmarshal(w.Body.Bytes(), &got)

			if got["code"] != tt.wantBody["code"] {
				t.Errorf("ChangePassword() code = %v, want %v", got["code"], tt.wantBody["code"])
			}
			if got["message"] != tt.wantBody["message"] {
				t.Errorf("ChangePassword() message = %v, want %v", got["message"], tt.wantBody["message"])
			}
		})
	}
}

func TestUserHandler_Logout(t *testing.T) {
	handler, router, mock := setupTestHandler()
	router.POST("/logout", handler.Logout)

	tests := []struct {
		name       string
		setupAuth  func(*http.Request)
		setupMock  func()
		wantStatus int
		wantBody   map[string]interface{}
	}{
		{
			name: "登出成功",
			setupAuth: func(req *http.Request) {
				req.Header.Set("Authorization", "Bearer valid-token")
			},
			setupMock: func() {
				mock.logoutErr = nil
			},
			wantStatus: http.StatusOK,
			wantBody: map[string]interface{}{
				"code":    float64(http.StatusOK),
				"message": "登出成功",
			},
		},
		{
			name: "未提供token",
			setupAuth: func(req *http.Request) {
				// 不设置token
			},
			setupMock:  func() {},
			wantStatus: http.StatusUnauthorized,
			wantBody: map[string]interface{}{
				"code":    float64(http.StatusUnauthorized),
				"message": "未授权",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setupMock != nil {
				tt.setupMock()
			}

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/logout", nil)

			if tt.setupAuth != nil {
				tt.setupAuth(req)
			}

			router.ServeHTTP(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("Logout() status = %v, want %v", w.Code, tt.wantStatus)
			}

			var got map[string]interface{}
			json.Unmarshal(w.Body.Bytes(), &got)

			if got["code"] != tt.wantBody["code"] {
				t.Errorf("Logout() code = %v, want %v", got["code"], tt.wantBody["code"])
			}
			if got["message"] != tt.wantBody["message"] {
				t.Errorf("Logout() message = %v, want %v", got["message"], tt.wantBody["message"])
			}
		})
	}
}
