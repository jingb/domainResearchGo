package model

// TotalTrafficAndEngagementResp 定义了流量分析的响应结构
type TotalTrafficAndEngagementResp struct {
	Meta struct {
		Request struct {
			Granularity    string      `json:"granularity"`
			MainDomainOnly bool        `json:"main_domain_only"`
			Mtd            bool        `json:"mtd"`
			ShowVerified   bool        `json:"show_verified"`
			State          interface{} `json:"state"`
			Format         string      `json:"format"`
			Domain         string      `json:"domain"`
			StartDate      string      `json:"start_date"`
			EndDate        string      `json:"end_date"`
			Country        string      `json:"country"`
		} `json:"request"`
		Status      string `json:"status"`
		LastUpdated string `json:"last_updated"`
	} `json:"meta"`
	Visits []struct {
		Date   string  `json:"date"`
		Visits float64 `json:"visits"`
	} `json:"visits"`
}
