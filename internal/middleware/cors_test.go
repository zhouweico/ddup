package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestCORSMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name            string
		method          string
		origin          string
		expectedHeaders map[string]string
		expectedStatus  int
	}{
		{
			name:   "正常 CORS 请求",
			method: "GET",
			origin: "http://localhost:3000",
			expectedHeaders: map[string]string{
				"Access-Control-Allow-Origin":      "http://localhost:3000",
				"Access-Control-Allow-Methods":     "GET,POST,PUT,PATCH,DELETE,OPTIONS",
				"Access-Control-Allow-Headers":     "Origin,Content-Type,Accept,Authorization",
				"Access-Control-Allow-Credentials": "true",
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:   "预检请求",
			method: "OPTIONS",
			origin: "http://localhost:3000",
			expectedHeaders: map[string]string{
				"Access-Control-Allow-Origin":      "http://localhost:3000",
				"Access-Control-Allow-Methods":     "GET,POST,PUT,PATCH,DELETE,OPTIONS",
				"Access-Control-Allow-Headers":     "Origin,Content-Type,Accept,Authorization",
				"Access-Control-Allow-Credentials": "true",
				"Access-Control-Max-Age":           "86400",
			},
			expectedStatus: http.StatusNoContent,
		},
		{
			name:   "无 Origin 请求",
			method: "GET",
			origin: "",
			expectedHeaders: map[string]string{
				"Access-Control-Allow-Origin": "*",
			},
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := gin.New()
			router.Use(Cors())

			router.GET("/test", func(c *gin.Context) {
				c.Status(http.StatusOK)
			})

			w := httptest.NewRecorder()
			req := httptest.NewRequest(tt.method, "/test", nil)
			if tt.origin != "" {
				req.Header.Set("Origin", tt.origin)
			}

			router.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("状态码 = %v, want %v", w.Code, tt.expectedStatus)
			}

			for key, expectedValue := range tt.expectedHeaders {
				if got := w.Header().Get(key); got != expectedValue {
					t.Errorf("Header[%s] = %v, want %v", key, got, expectedValue)
				}
			}
		})
	}
}
