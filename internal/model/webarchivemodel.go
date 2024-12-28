package model

import "time"

// WebArchiveResponse 定义了域名收录信息的响应结构
type WebArchiveResponse struct {
	CreateTime time.Time `json:"create_time"`
	Original   string    `json:"original"`
}
