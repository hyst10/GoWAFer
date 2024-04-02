package middleware

import (
	"GoWAFer/pkg/utils/api_helper"
	"github.com/gin-gonic/gin"
	"regexp"
)

// XSSDetectMiddleware 检测xss攻击的中间件
func XSSDetectMiddleware(rules []*regexp.Regexp) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 检查是否为白名单IP
		if skip, _ := c.Get("isWhiteListIP"); skip == true {
			c.Next()
			return
		}

		// 检查是否为白名单路由
		if skip, _ := c.Get("isWhiteRoute"); skip == true {
			c.Next()
			return
		}

		// 检查URL查询参数
		query := c.Request.URL.Query()
		for _, values := range query {
			for _, value := range values {
				for _, pattern := range rules {
					if pattern.MatchString(value) {
						c.Set("BlockedBy", "xss攻击防护中间件")
						c.Set("BlockReason", "查询参数中检测到xss攻击")
						api_helper.ForbiddenHandler(c, "检测到xss攻击，禁止访问！")
						c.Abort()
						return
					}
				}
			}
		}

		// 检查表单数据
		contentType := c.GetHeader("Content-Type")
		if contentType == "application/x-www-form-urlencoded" || contentType == "multipart/form-data" {
			if err := c.Request.ParseForm(); err == nil {
				for _, values := range c.Request.PostForm {
					for _, value := range values {
						for _, pattern := range rules {
							if pattern.MatchString(value) {
								c.Set("BlockedBy", "xss攻击防护中间件")
								c.Set("BlockReason", "表单数据中检测到xss攻击")
								api_helper.ForbiddenHandler(c, "检测到xss攻击，禁止访问！")
								c.Abort()
								return
							}
						}
					}
				}
			}
		}

		c.Next()
	}
}
