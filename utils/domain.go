package utils

import (
	"golang.org/x/net/publicsuffix"
	"net"
	"strings"
)

// TopLevelDomain 获取顶级域名
func TopLevelDomain(domain string) (string, error) {
	return publicsuffix.EffectiveTLDPlusOne(domain)
}

func GetHostTLD(host string) (string, error) {
	h, _, err := net.SplitHostPort(host)
	if err != nil {
		h = host // 如果没有端口号，直接使用 host
	}
	if strings.HasPrefix(h, "[") && strings.Contains(host, "]") {
		// IPv6 地址
		end := strings.Index(h, "]")
		h = h[:end+1]
	}
	return TopLevelDomain(h)
}
