package errors

import (
	"errors"
	"net/http"
	"testing"
)

func TestAppError_Error(t *testing.T) {
	tests := []struct {
		name string
		err  *AppError
		want string
	}{
		{
			name: "仅包含消息的错误",
			err: &AppError{
				Message: "测试错误",
			},
			want: "测试错误",
		},
		{
			name: "包含原始错误的错误",
			err: &AppError{
				Message: "测试错误",
				Err:     errors.New("原始错误"),
			},
			want: "测试错误: 原始错误",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.err.Error(); got != tt.want {
				t.Errorf("AppError.Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNew(t *testing.T) {
	originalErr := errors.New("原始错误")
	appErr := New(http.StatusBadRequest, "测试错误", originalErr)

	if appErr.Code != http.StatusBadRequest {
		t.Errorf("New().Code = %v, want %v", appErr.Code, http.StatusBadRequest)
	}
	if appErr.Message != "测试错误" {
		t.Errorf("New().Message = %v, want %v", appErr.Message, "测试错误")
	}
	if appErr.Err != originalErr {
		t.Errorf("New().Err = %v, want %v", appErr.Err, originalErr)
	}
}

func TestWrap(t *testing.T) {
	t.Run("包装普通错误", func(t *testing.T) {
		err := errors.New("原始错误")
		wrapped := Wrap(err, "包装错误")

		if wrapped.Code != http.StatusInternalServerError {
			t.Errorf("Wrap().Code = %v, want %v", wrapped.Code, http.StatusInternalServerError)
		}
		if wrapped.Message != "包装错误" {
			t.Errorf("Wrap().Message = %v, want %v", wrapped.Message, "包装错误")
		}
		if wrapped.Err != err {
			t.Errorf("Wrap().Err = %v, want %v", wrapped.Err, err)
		}
	})

	t.Run("包装 AppError", func(t *testing.T) {
		original := &AppError{
			Code:    http.StatusBadRequest,
			Message: "原始错误",
		}
		wrapped := Wrap(original, "新消息")

		if wrapped != original {
			t.Error("Wrap() 应该返回原始 AppError")
		}
	})
}
