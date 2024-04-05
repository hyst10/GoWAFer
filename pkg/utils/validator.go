package utils

import (
	"net"
	"regexp"
)

// ValidateRouting 检验是否为路由
func ValidateRouting(routing string) bool {
	pattern := `^/[A-Za-z0-9/_]*$`
	matched, _ := regexp.MatchString(pattern, routing)
	if !matched {
		return false
	}
	return true
}

// ValidateIP 检验是否为正确的IP
func ValidateIP(ip string) bool {
	if net.ParseIP(ip) == nil {
		return false
	}
	return true
}
