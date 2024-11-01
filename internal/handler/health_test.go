package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

// 添加 Mock 数据库连接
type mockDB struct {
	pingErr error
}

func (m *mockDB) Ping() error {
	return m.pingErr
}

// 修改为接口类型
type DBInterface interface {
	Ping() error
}

// 在测试中使用接口
var DB DBInterface

func TestHealthHandler_Check(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name       string
		dbPingErr  error
		wantStatus int
		wantBody   map[string]interface{}
	}{
		{
			name:       "服务正常",
			dbPingErr:  nil,
			wantStatus: http.StatusOK,
			wantBody: map[string]interface{}{
				"code":    float64(http.StatusOK),
				"message": "服务正常",
			},
		},
		{
			name:       "数据库异常",
			dbPingErr:  errors.New("数据库连接失败"),
			wantStatus: http.StatusServiceUnavailable,
			wantBody: map[string]interface{}{
				"code":    float64(http.StatusServiceUnavailable),
				"message": "数据库连接异常",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 使用接口而不是直接赋值给 db.DB
			DB = &mockDB{pingErr: tt.dbPingErr}

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
