package middleware

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

func TestLoggerMiddleware(t *testing.T) {
	// 设置测试模式
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		method         string
		path           string
		body           string
		setupRequest   func(*http.Request)
		expectedStatus int
		expectedLog    string
	}{
		{
			name:   "正常请求",
			method: "GET",
			path:   "/test",
			setupRequest: func(req *http.Request) {
				req.Header.Set("User-Agent", "test-agent")
			},
			expectedStatus: http.StatusOK,
			expectedLog:    "GET /test",
		},
		{
			name:   "带请求体的POST请求",
			method: "POST",
			path:   "/test",
			body:   `{"key":"value"}`,
			setupRequest: func(req *http.Request) {
				req.Header.Set("Content-Type", "application/json")
			},
			expectedStatus: http.StatusOK,
			expectedLog:    "POST /test",
		},
		{
			name:   "错误请求",
			method: "GET",
			path:   "/error",
			setupRequest: func(req *http.Request) {
				req.Header.Set("X-Request-ID", "test-id")
			},
			expectedStatus: http.StatusInternalServerError,
			expectedLog:    "GET /error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 创建测试路由
			r := gin.New()
			r.Use(Logger())

			// 添加测试路由
			r.Handle(tt.method, "/test", func(c *gin.Context) {
				c.Status(http.StatusOK)
			})
			r.Handle(tt.method, "/error", func(c *gin.Context) {
				c.Status(http.StatusInternalServerError)
			})

			// 创建测试请求
			var reqBody *bytes.Buffer
			if tt.body != "" {
				reqBody = bytes.NewBufferString(tt.body)
			} else {
				reqBody = bytes.NewBuffer(nil)
			}

			w := httptest.NewRecorder()
			req := httptest.NewRequest(tt.method, tt.path, reqBody)

			if tt.setupRequest != nil {
				tt.setupRequest(req)
			}

			// 执行请求
			r.ServeHTTP(w, req)

			// 验证响应状态码
			if w.Code != tt.expectedStatus {
				t.Errorf("状态码 = %v, want %v", w.Code, tt.expectedStatus)
			}

			// 等待日志写入
			time.Sleep(100 * time.Millisecond)

			// 验证日志内容（这里需要根据实际的日志实现来调整）
			//
		})
	}
}
