package utils

import (
	"errors"
	"golang.org/x/net/publicsuffix"
	"net"
	"strings"
)

// TopLevelDomain 获取顶级域名
func TopLevelDomain(domain string) (string, error) {
	return publicsuffix.EffectiveTLDPlusOne(domain)
}

func GetHostAddr(host string) (string, error) {
	h, _, err := net.SplitHostPort(host)
	if err != nil {
		var e *net.AddrError
		if errors.As(err, &e) {
			h = e.Addr
		} else {
			return "", err
		}
	}
	if strings.HasPrefix(h, "[") && strings.Contains(host, "]") {
		// IPv6 地址
		end := strings.Index(h, "]")
		h = h[:end+1]
	}
	return h, nil
}

func GetHostTLD(host string) (string, error) {
	h, err := GetHostAddr(host)
	if err != nil {
		return "", err
	}
	return TopLevelDomain(h)
}
