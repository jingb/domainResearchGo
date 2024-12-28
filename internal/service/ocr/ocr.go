package ocr

import (
	"context"
	"net/url"
)

// OCRService 定义OCR服务的接口
type OCRService interface {
	// RecognizeDomains 从图片中识别并提取合法域名
	// 返回经过验证的合法域名URL对象列表
	RecognizeDomains(ctx context.Context, imageBytes []byte) ([]*url.URL, error)
}

type OCRResponse struct {
	Texts []string `json:"texts"`
}
