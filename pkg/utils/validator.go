package utils

import (
	"regexp"
)

// ValidatorRouting 检验是否为标准的路由
func ValidatorRouting(routing string) bool {
	pattern := `^/[A-Za-z0-9/_]*$`
	matched, _ := regexp.MatchString(pattern, routing)
	if !matched {
		return false
	}
	return true
}
