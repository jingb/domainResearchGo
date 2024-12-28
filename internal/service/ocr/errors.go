package ocr

import "errors"

// OCR服务内部错误定义
var (
	ErrEmptyImage   = errors.New("empty image")
	ErrInvalidImage = errors.New("invalid image format")
	ErrServiceError = errors.New("ocr service error")
)
