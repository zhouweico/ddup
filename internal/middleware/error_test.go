package middleware

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"ddup-apis/internal/errors"

	"github.com/gin-gonic/gin"
)

func TestErrorHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name       string
		setup      func(*gin.Engine)
		wantStatus int
		wantBody   map[string]interface{}
	}{
		{
			name: "处理AppError",
			setup: func(r *gin.Engine) {
				r.GET("/test", func(c *gin.Context) {
					appErr := errors.New(http.StatusBadRequest, "测试错误", nil)
					_ = c.Error(appErr)
				})
			},
			wantStatus: http.StatusBadRequest,
			wantBody: map[string]interface{}{
				"code":    float64(http.StatusBadRequest),
				"message": "测试错误",
			},
		},
		{
			name: "处理普通错误",
			setup: func(r *gin.Engine) {
				r.GET("/test", func(c *gin.Context) {
					err := errors.New(http.StatusInternalServerError, "系统错误", nil)
					_ = c.Error(err)
				})
			},
			wantStatus: http.StatusInternalServerError,
			wantBody: map[string]interface{}{
				"code":    float64(http.StatusInternalServerError),
				"message": "系统错误",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := gin.New()
			r.Use(ErrorHandler())
			tt.setup(r)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/test", nil)
			r.ServeHTTP(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("ErrorHandler() status = %v, want %v", w.Code, tt.wantStatus)
			}

			var got map[string]interface{}
			if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
				t.Fatalf("解析响应体失败: %v", err)
			}

			if got["code"] != tt.wantBody["code"] {
				t.Errorf("ErrorHandler() code = %v, want %v", got["code"], tt.wantBody["code"])
			}
			if got["message"] != tt.wantBody["message"] {
				t.Errorf("ErrorHandler() message = %v, want %v", got["message"], tt.wantBody["message"])
			}
		})
	}
}
