package errors

import (
	"fmt"
	"net/http"
)

// AppError 自定义错误类型
type AppError struct {
	Code    int    // HTTP 状态码
	Message string // 错误信息
	Err     error  // 原始错误
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

// 预定义错误
var (
	ErrInvalidRequest = &AppError{
		Code:    http.StatusBadRequest,
		Message: "无效的请求参数",
	}

	ErrUnauthorized = &AppError{
		Code:    http.StatusUnauthorized,
		Message: "未授权访问",
	}

	ErrForbidden = &AppError{
		Code:    http.StatusForbidden,
		Message: "禁止访问",
	}

	ErrNotFound = &AppError{
		Code:    http.StatusNotFound,
		Message: "资源不存在",
	}

	ErrInternalServer = &AppError{
		Code:    http.StatusInternalServerError,
		Message: "服务器内部错误",
	}
)

// New 创建新的应用错误
func New(code int, message string, err error) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

// Wrap 包装已有错误
func Wrap(err error, message string) *AppError {
	if appErr, ok := err.(*AppError); ok {
		return appErr
	}
	return &AppError{
		Code:    http.StatusInternalServerError,
		Message: message,
		Err:     err,
	}
}
