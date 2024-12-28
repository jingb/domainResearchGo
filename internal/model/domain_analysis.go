package model

// DomainAnalysis 表示单个域名的分析结果
type DomainAnalysis struct {
	Domain                        string                        `json:"domain"`
	WebArchiveResponse            WebArchiveResponse            `json:"web_archive_response"`
	TotalTrafficAndEngagementResp TotalTrafficAndEngagementResp `json:"total_traffic_and_engagement_response"`
}
