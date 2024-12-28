package errors

import (
	"errors"
	"fmt"
)

// ErrorType 错误类型
type ErrorType uint

const (
	// ErrorTypeClient 客户端错误 (4xx)
	ErrorTypeClient ErrorType = iota
	// ErrorTypeServer 服务端错误 (5xx)
	ErrorTypeServer
)

// Error 自定义错误结构
type Error struct {
	Type    ErrorType // 错误类型
	Message string    // 错误信息
	Err     error     // 原始错误
}

// Error 实现error接口
func (e *Error) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

// Unwrap 支持错误链
func (e *Error) Unwrap() error {
	return e.Err
}

// NewClientError 创建客户端错误
func NewClientError(message string, err error) *Error {
	return &Error{
		Type:    ErrorTypeClient,
		Message: message,
		Err:     err,
	}
}

// NewServerError 创建服务端错误
func NewServerError(message string, err error) *Error {
	return &Error{
		Type:    ErrorTypeServer,
		Message: message,
		Err:     err,
	}
}

// IsClientError 判断是否为客户端错误
func IsClientError(err error) bool {
	var e *Error
	if err == nil {
		return false
	}
	if As(err, &e) {
		return e.Type == ErrorTypeClient
	}
	return false
}

// IsServerError 判断是否为服务端错误
func IsServerError(err error) bool {
	var e *Error
	if err == nil {
		return false
	}
	if As(err, &e) {
		return e.Type == ErrorTypeServer
	}
	return false
}

// As 包装标准库的errors.As
func As(err error, target interface{}) bool {
	return errors.As(err, target)
}
