package handler

import (
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

func TestSendSuccess(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	tests := []struct {
		name    string
		message string
		data    interface{}
		want    Response
	}{
		{
			name:    "成功响应-带数据",
			message: "操作成功",
			data:    map[string]string{"key": "value"},
			want: Response{
				Code:    200,
				Message: "操作成功",
				Data:    map[string]string{"key": "value"},
			},
		},
		{
			name:    "成功响应-无数据",
			message: "操作成功",
			data:    nil,
			want: Response{
				Code:    200,
				Message: "操作成功",
				Data:    nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SendSuccess(c, tt.message, tt.data)
			if w.Code != 200 {
				t.Errorf("SendSuccess() status = %v, want %v", w.Code, 200)
			}
		})
	}
}

func TestSendError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	tests := []struct {
		name       string
		status     int
		message    string
		wantStatus int
	}{
		{
			name:       "错误响应-400",
			status:     400,
			message:    "请求参数错误",
			wantStatus: 400,
		},
		{
			name:       "错误响应-500",
			status:     500,
			message:    "服务器内部错误",
			wantStatus: 500,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SendError(c, tt.status, tt.message)
			if w.Code != tt.wantStatus {
				t.Errorf("SendError() status = %v, want %v", w.Code, tt.wantStatus)
			}
		})
	}
}

func TestTokenInfo(t *testing.T) {
	now := time.Now()
	expiresIn := int64(3600)
	expiredAt := now.Add(time.Duration(expiresIn) * time.Second)

	info := TokenInfo{
		Token:     "test-token",
		CreatedAt: now,
		ExpiresIn: expiresIn,
		ExpiredAt: expiredAt,
	}

	if info.Token != "test-token" {
		t.Errorf("TokenInfo.Token = %v, want %v", info.Token, "test-token")
	}
	if info.ExpiresIn != expiresIn {
		t.Errorf("TokenInfo.ExpiresIn = %v, want %v", info.ExpiresIn, expiresIn)
	}
	if !info.CreatedAt.Equal(now) {
		t.Errorf("TokenInfo.CreatedAt = %v, want %v", info.CreatedAt, now)
	}
}
