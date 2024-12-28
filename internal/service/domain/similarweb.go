package domain

import (
	"context"
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
	similarWebBaseURL = "https://api.similarweb.com/v1/website"
)

var (
	similarWebInstance SimilarWeb
	similarWebOnce     sync.Once
)

// TrafficQuery 定义流量查询的参数
type TrafficQuery struct {
	// 数据粒度: "daily", "weekly", "monthly"
	Granularity string `json:"granularity"`
	// 是否只查询主域名
	MainDomainOnly bool `json:"main_domain_only"`
	// 是否包含当月至今的数据
	MonthToDate bool `json:"mtd"`
	// 是否只显示已验证的数据
	ShowVerified bool `json:"show_verified"`
	// 返回格式: "json"
	Format string `json:"format"`
	// 开始日期 "YYYY-MM"
	StartDate string `json:"start_date"`
	// 结束日期 "YYYY-MM"
	EndDate string `json:"end_date"`
	// 国家代码 "world" 或 ISO 2字母代码
	Country string `json:"country"`
}

type SimilarWeb interface {
	// TotalTrafficAndEngagement 获取某个域名的总流量
	// https://developers.similarweb.com/reference/visits
	TotalTrafficAndEngagement(ctx context.Context, query TrafficQuery, domain *url.URL) (model.TotalTrafficAndEngagementResp, error)
}

// similarWebClient 是 SimilarWeb API 的客户端
type similarWebClient struct {
	client  *http.Client
	apiKey  string
	baseURL string
}

// Config SimilarWeb客户端配置
type SimilarWebConfig struct {
	APIKey string
}

// newSimilarWebClient 创建一个新的 SimilarWeb 客���端
func newSimilarWebClient(config *SimilarWebConfig) SimilarWeb {
	return &similarWebClient{
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
		apiKey:  config.APIKey,
		baseURL: similarWebBaseURL,
	}
}

// GetSimilarWeb 返回 SimilarWeb 的全局单例实例
func GetSimilarWeb(config *SimilarWebConfig) SimilarWeb {
	similarWebOnce.Do(func() {
		similarWebInstance = newSimilarWebClient(config)
	})
	return similarWebInstance
}

// TotalTrafficAndEngagement 实现 SimilarWeb 接口
func (c *similarWebClient) TotalTrafficAndEngagement(ctx context.Context, query TrafficQuery, domain *url.URL) (model.TotalTrafficAndEngagementResp, error) {
	// 构建请求URL
	endpointURL := fmt.Sprintf("%s/%s/total-traffic-and-engagement/visits", c.baseURL, domain.Host)

	// 构建查询参数
	params := url.Values{}
	if query.Granularity != "" {
		params.Add("granularity", query.Granularity)
	}
	if query.MainDomainOnly {
		params.Add("main_domain_only", "true")
	}
	if query.MonthToDate {
		params.Add("mtd", "true")
	}
	if query.ShowVerified {
		params.Add("show_verified", "true")
	}
	if query.Format != "" {
		params.Add("format", query.Format)
	}
	if query.StartDate != "" {
		params.Add("start_date", query.StartDate)
	}
	if query.EndDate != "" {
		params.Add("end_date", query.EndDate)
	}
	if query.Country != "" {
		params.Add("country", query.Country)
	}

	// 添加查询参数到URL
	requestURL := endpointURL
	if len(params) > 0 {
		requestURL += "?" + params.Encode()
	}

	// 创建请求
	req, err := http.NewRequestWithContext(ctx, "GET", requestURL, nil)
	if err != nil {
		return model.TotalTrafficAndEngagementResp{}, fmt.Errorf("create request failed: %v (URL: %s)", err, requestURL)
	}

	// 添加认证头
	req.Header.Set("api_key", c.apiKey)

	// 发送请求
	resp, err := c.client.Do(req)
	if err != nil {
		return model.TotalTrafficAndEngagementResp{}, fmt.Errorf("request failed: %v (URL: %s)", err, requestURL)
	}
	defer resp.Body.Close()

	// 读取响应体
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return model.TotalTrafficAndEngagementResp{}, fmt.Errorf("read response failed: %v (URL: %s, StatusCode: %d)",
			err, requestURL, resp.StatusCode)
	}

	// 检查响应状态
	if resp.StatusCode != http.StatusOK {
		return model.TotalTrafficAndEngagementResp{}, fmt.Errorf("API request failed with status %d (URL: %s, Response: %s)",
			resp.StatusCode, requestURL, string(body))
	}

	// 解析响应
	var result model.TotalTrafficAndEngagementResp
	if err := json.Unmarshal(body, &result); err != nil {
		return model.TotalTrafficAndEngagementResp{}, fmt.Errorf("parse response failed: %v (URL: %s, Response: %s)",
			err, requestURL, string(body))
	}

	return result, nil
}
