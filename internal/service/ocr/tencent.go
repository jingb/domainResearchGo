package ocr

import (
	"context"
	"domain-analyzer/config"
	"domain-analyzer/internal/pkg/domainutil"
	"encoding/base64"
	"net/url"

	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	ocr "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ocr/v20181119"
)

// TencentOCR 腾讯云OCR服务实现
type TencentOCR struct {
	client *ocr.Client
}

// NewTencentOCR 创建新的腾讯云OCR服务实例
func NewTencentOCR(config *config.Config) (*TencentOCR, error) {
	credential := common.NewCredential(
		config.TencentCloud.SecretId,
		config.TencentCloud.SecretKey,
	)

	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = "ocr.tencentcloudapi.com"

	client, err := ocr.NewClient(credential, config.TencentCloud.Region, cpf)
	if err != nil {
		return nil, err
	}

	return &TencentOCR{
		client: client,
	}, nil
}

// RecognizeDomains 实现OCRService接口，从图片中识别并提取域名
func (t *TencentOCR) RecognizeDomains(ctx context.Context, imageBytes []byte) ([]*url.URL, error) {
	// 将图片转换为Base64
	base64Img := base64.StdEncoding.EncodeToString(imageBytes)

	// 创建通用OCR请求
	request := ocr.NewGeneralBasicOCRRequest()
	request.ImageBase64 = common.StringPtr(base64Img)

	// 调用OCR API
	response, err := t.client.GeneralBasicOCR(request)
	if err != nil {
		return nil, err
	}

	// 提取所有识别出的文本
	var texts []string
	for _, textDetection := range response.Response.TextDetections {
		if textDetection.DetectedText != nil {
			texts = append(texts, *textDetection.DetectedText)
		}
	}

	// 从文本中提取域名
	return domainutil.ExtractDomains(texts)
}
