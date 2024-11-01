package router

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestSetupRouter(t *testing.T) {
	// 设置测试模式
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name         string
		path         string
		method       string
		wantStatus   int
		setupHeaders func(*http.Request)
	}{
		{
			name:       "健康检查路由",
			path:       "/health",
			method:     "GET",
			wantStatus: http.StatusOK,
		},
		{
			name:       "Swagger文档路由",
			path:       "/swagger/index.html",
			method:     "GET",
			wantStatus: http.StatusOK,
		},
		{
			name:       "注册路由",
			path:       "/api/v1/register",
			method:     "POST",
			wantStatus: http.StatusBadRequest, // 因为没有请求体
		},
		{
			name:       "登录路由",
			path:       "/api/v1/login",
			method:     "POST",
			wantStatus: http.StatusBadRequest, // 因为没有请求体
		},
		{
			name:   "需要认证的路由",
			path:   "/api/v1/user",
			method: "GET",
			setupHeaders: func(r *http.Request) {
				r.Header.Set("Authorization", "Bearer invalid-token")
			},
			wantStatus: http.StatusUnauthorized,
		},
	}

	// 初始化路由
	r := SetupRouter()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req := httptest.NewRequest(tt.method, tt.path, nil)

			if tt.setupHeaders != nil {
				tt.setupHeaders(req)
			}

			r.ServeHTTP(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("路由 %s 返回状态码 %d，期望 %d", tt.name, w.Code, tt.wantStatus)
			}
		})
	}
}
