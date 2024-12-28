package handler

import (
	"bytes"
	"context"
	"domain-analyzer/internal/model"
	"domain-analyzer/internal/pkg/errors"
	"domain-analyzer/internal/service/domain"
	"domain-analyzer/internal/service/ocr"
	"io"
	"net/http"
)

type UploadHandler struct {
	ocrService ocr.OCRService
	webArchive domain.WebArchive
}

func NewUploadHandler(ocrService ocr.OCRService) Handler {
	return &UploadHandler{
		ocrService: ocrService,
		webArchive: domain.GetWebArchive(),
	}
}

// UploadResponse 表示图片上传并分析后的响应结果
type UploadResponse struct {
	Domains []model.DomainAnalysis `json:"domains"`
}

func (h *UploadHandler) Handle(ctx context.Context, req *http.Request) (interface{}, error) {
	// 从multipart form中获取文件
	if err := req.ParseMultipartForm(32 << 20); err != nil {
		return nil, errors.NewClientError("解析上传请求失败", err)
	}

	file, _, err := req.FormFile("image")
	if err != nil {
		return nil, errors.NewClientError("未找到上传的图片文件", err)
	}
	defer file.Close()

	// 读取图片数据
	var buf bytes.Buffer
	if _, err = io.Copy(&buf, file); err != nil {
		return nil, errors.NewClientError("图片文件读取失败", err)
	}

	// 调用OCR服务识别域名
	domains, err := h.ocrService.RecognizeDomains(ctx, buf.Bytes())
	if err != nil {
		return nil, err
	}

	ret := &UploadResponse{}
	for _, domain := range domains {
		analysisResult, err := h.webArchive.RecognizeDomains(ctx, domain)
		if err != nil {
			return nil, err
		}
		ret.Domains = append(ret.Domains, model.DomainAnalysis{
			Domain:             domain.Host,
			WebArchiveResponse: analysisResult,
		})
	}
	return ret, nil
}
