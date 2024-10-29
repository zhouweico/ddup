package utils

import (
	"fmt"
)

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func SuccessResponse(data interface{}) Response {
	return Response{
		Code:    200,
		Message: "success",
		Data:    data,
	}
}

func ErrorResponse(err interface{}) Response {
	return Response{
		Code:    400,
		Message: fmt.Sprint(err),
	}
}
