package domainutil

import (
	"net/url"
	"regexp"
	"strings"
)

var (
	// 域名的基本正则表达式
	// 1. 允许字母、数字、连字符
	// 2. 必须包含至少一个点号
	// 3. 不允许连续的点号或连字符
	// 4. 不允许开头或结尾是连字符
	domainRegex = regexp.MustCompile(`^([a-zA-Z0-9]([a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?\.)+[a-zA-Z]{2,}$`)

	// 常见的顶级域名后缀
	commonTLDs = map[string]bool{
		"com": true, "org": true, "net": true, "edu": true, "gov": true,
		"ru": true, "cn": true, "uk": true, "jp": true, "de": true,
		"fr": true, "br": true, "it": true, "pl": true, "in": true,
		"info": true, "biz": true, "io": true, "co": true, "me": true,
	}
)

// ExtractDomains 从文本列表中提取合法的域名
func ExtractDomains(texts []string) ([]*url.URL, error) {
	var results []*url.URL
	seen := make(map[string]bool) // 用于去重

	for _, text := range texts {
		// 预处理文本
		text = strings.TrimSpace(text)
		text = strings.ToLower(text)

		// 如果文本不包含点号，跳过
		if !strings.Contains(text, ".") {
			continue
		}

		// 尝试解析为URL
		u, err := url.Parse("http://" + text)
		if err != nil {
			continue
		}

		// 获取主机名部分
		host := u.Hostname()

		// 如果已经处理过这个域名，跳过
		if seen[host] {
			continue
		}

		// 验证域名格式
		if !domainRegex.MatchString(host) {
			continue
		}

		// 验证顶级域名
		parts := strings.Split(host, ".")
		if len(parts) < 2 {
			continue
		}
		tld := parts[len(parts)-1]
		if !commonTLDs[tld] {
			continue
		}

		// 域名长度检查
		if len(host) < 3 || len(host) > 255 {
			continue
		}

		// 标记为已处理
		seen[host] = true

		// 创建新的URL对象（只保留域名部分）
		domainURL, _ := url.Parse("http://" + host)
		results = append(results, domainURL)
	}

	return results, nil
}
