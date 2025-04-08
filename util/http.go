package util

import (
	"net"
	"net/http"
	"strings"
)

// GetHttpClientIP 获取本机IP地址
// 后续可以通过 net.ParseIP 验证IP地址的有效性
func GetHttpClientIP(r http.Request) string {
	// 首先检查X-Forwarded-For
	forwarded := r.Header.Get("X-Forwarded-For")
	if forwarded != "" {
		// 取第一个IP
		ips := strings.Split(forwarded, ",")
		ip := strings.TrimSpace(ips[0])
		if net.ParseIP(ip) != nil {
			return ip
		}
	}

	// 然后检查X-Real-IP
	realIP := r.Header.Get("X-Real-IP")
	if realIP != "" && net.ParseIP(realIP) != nil {
		return realIP
	}

	// 最后使用RemoteAddr
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return ip
}
