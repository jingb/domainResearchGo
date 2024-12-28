package domain

import (
	"context"
	"domain-analyzer/config"
	"domain-analyzer/internal/model"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sync"
	"time"
)

const (
	webArchiveBaseURL = "https://web.archive.org/cdx/search/cdx"
)

var (
	instance WebArchive
	once     sync.Once
	cfg      *config.Config
)

// WebArchive 定义了与 Internet Archive 交互的接口
type WebArchive interface {
	// RecognizeDomains 从 Web Archive 服务获取某个域名最早的收录时间
	RecognizeDomains(ctx context.Context, domain *url.URL) (model.WebArchiveResponse, error)
}

// cdxClient 是 Web Archive CDX Server API 的客户端
type cdxClient struct {
	client *http.Client
}

// newCDXClient 创建一个新的 CDX 客户端
func newCDXClient(proxyURL string) WebArchive {
	transport := &http.Transport{}

	// 如果提供了代理URL，则配置代理
	if proxyURL != "" {
		if proxy, err := url.Parse(proxyURL); err == nil {
			transport.Proxy = http.ProxyURL(proxy)
		}
	}

	return &cdxClient{
		client: &http.Client{
			Transport: transport,
			Timeout:   10 * time.Second,
		},
	}
}

// 添加一个初始化配置的函数
func InitWebArchive(config *config.Config) {
	cfg = config
}

// GetWebArchive 返回 WebArchive 的全局单例实例
func GetWebArchive() WebArchive {
	once.Do(func() {
		// 使用包级变量中的配置
		proxyURL := ""
		if cfg != nil {
			proxyURL = cfg.WebArchive.ProxyURL
		}
		instance = newCDXClient(proxyURL)
	})
	return instance
}

// CDXResponse CDX API 的响应格式
// 第一行是字段名，第二行是实际数据
// 例如: [["timestamp","original"],["20100615142933","http://example.com"]]
type CDXResponse [][]string

// RecognizeDomains 实现 WebArchive 接口
func (c *cdxClient) RecognizeDomains(ctx context.Context, domain *url.URL) (model.WebArchiveResponse, error) {
	// 构建查询 URL
	params := url.Values{}
	params.Add("url", domain.Host)
	params.Add("output", "json")
	params.Add("fl", "timestamp,original")
	params.Add("sort", "timestamp")
	params.Add("limit", "1")

	queryURL := fmt.Sprintf("%s?%s", webArchiveBaseURL, params.Encode())

	// 创建请求
	req, err := http.NewRequestWithContext(ctx, "GET", queryURL, nil)
	if err != nil {
		return model.WebArchiveResponse{}, fmt.Errorf("create request failed: %w", err)
	}

	// 发送请求
	resp, err := c.client.Do(req)
	if err != nil {
		return model.WebArchiveResponse{}, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return model.WebArchiveResponse{}, fmt.Errorf("read response failed: %w", err)
	}

	// 解析响应
	var cdxResp CDXResponse
	if err := json.Unmarshal(body, &cdxResp); err != nil {
		return model.WebArchiveResponse{}, fmt.Errorf("parse response failed: %w", err)
	}

	// 检查响应格式
	if len(cdxResp) < 2 {
		return model.WebArchiveResponse{}, fmt.Errorf("no archive found for domain: %s", domain.Host)
	}

	// 解析时间戳
	timestamp := cdxResp[1][0] // 获取第一条记录的时间戳
	original := cdxResp[1][1]
	createTime, err := time.Parse("20060102150405", timestamp)
	if err != nil {
		return model.WebArchiveResponse{}, fmt.Errorf("parse timestamp failed: %w", err)
	}

	return model.WebArchiveResponse{
		CreateTime: createTime,
		Original:   original,
	}, nil
}
