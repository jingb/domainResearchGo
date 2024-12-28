package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	TencentCloud struct {
		SecretId  string `json:"secret_id"`
		SecretKey string `json:"secret_key"`
		Region    string `json:"region"`
	} `json:"tencent_cloud"`
	Server struct {
		Port string `json:"port"`
	} `json:"server"`
	Analysis struct {
		TrafficThreshold int64 `json:"traffic_threshold"`
		DaysThreshold    int   `json:"days_threshold"`
	} `json:"analysis"`
	WebArchive struct {
		ProxyURL string `json:"proxy_url"`
	} `json:"web_archive"`
}

func LoadConfig(path string) (*Config, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config Config
	if err := json.Unmarshal(file, &config); err != nil {
		return nil, err
	}

	return &config, nil
}
