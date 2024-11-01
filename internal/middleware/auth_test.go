package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"ddup-apis/internal/service"
	"ddup-apis/internal/utils"

	"github.com/gin-gonic/gin"
)

func setupTestRouter(userService *service.UserService) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(JWTAuth(userService))
	return r
}

func TestJWTAuth(t *testing.T) {
	// 创建 mock 用户服务
	mockUserService := &service.UserService{}

	tests := []struct {
		name       string
		setupAuth  func() string
		wantStatus int
	}{
		{
			name: "无token",
			setupAuth: func() string {
				return ""
			},
			wantStatus: http.StatusUnauthorized,
		},
		{
			name: "无效token",
			setupAuth: func() string {
				return "invalid-token"
			},
			wantStatus: http.StatusUnauthorized,
		},
		{
			name: "有效token",
			setupAuth: func() string {
				token, _, _, _, _ := utils.GenerateToken(1, "testuser")
				return token
			},
			wantStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := setupTestRouter(mockUserService)
			r.GET("/test", func(c *gin.Context) {
				c.Status(http.StatusOK)
			})

			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/test", nil)

			if token := tt.setupAuth(); token != "" {
				req.Header.Set("Authorization", "Bearer "+token)
			}

			r.ServeHTTP(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("JWTAuth() status = %v, want %v", w.Code, tt.wantStatus)
			}
		})
	}
}
