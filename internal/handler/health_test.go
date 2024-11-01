package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"ddup-apis/internal/db"

	"github.com/gin-gonic/gin"
)

func TestHealthHandler_Check(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name       string
		setupDB    func()
		wantStatus int
		wantBody   map[string]interface{}
	}{
		{
			name: "数据库正常",
			setupDB: func() {
				// 模拟数据库正常
				db.DB = nil // 这里需要根据实际情况设置测试数据库
			},
			wantStatus: http.StatusOK,
			wantBody: map[string]interface{}{
				"code":    float64(http.StatusOK),
				"message": "服务正常",
				"data":    nil,
			},
		},
		{
			name: "数据库异常",
			setupDB: func() {
				// 模拟数据库异常
				db.DB = nil
			},
			wantStatus: http.StatusServiceUnavailable,
			wantBody: map[string]interface{}{
				"code":    float64(http.StatusServiceUnavailable),
				"message": "数据库连接异常",
			},
		},
		{
			name: "数据库连接超时",
			setupDB: func() {
				// 模拟数据库连接超时
				db.DB = nil // 设置一个会导致超时的测试数据库连接
			},
			wantStatus: http.StatusServiceUnavailable,
			wantBody: map[string]interface{}{
				"code":    float64(http.StatusServiceUnavailable),
				"message": "数据库连接异常",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupDB()

			r := gin.New()
			handler := NewHealthHandler()
			r.GET("/health", handler.Check)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/health", nil)
			r.ServeHTTP(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("HealthHandler.Check() status = %v, want %v", w.Code, tt.wantStatus)
			}

			var got map[string]interface{}
			if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
				t.Fatalf("解析响应体失败: %v", err)
			}

			if got["code"] != tt.wantBody["code"] {
				t.Errorf("HealthHandler.Check() code = %v, want %v", got["code"], tt.wantBody["code"])
			}
			if got["message"] != tt.wantBody["message"] {
				t.Errorf("HealthHandler.Check() message = %v, want %v", got["message"], tt.wantBody["message"])
			}
		})
	}
}
