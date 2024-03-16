package middleware

import (
	"GoWAFer/pkg/utils/api_handler"
	"github.com/gin-gonic/gin"
	"regexp"
)

// SqlInjectMiddleware sql注入检测中间件
func SqlInjectMiddleware(rules []*regexp.Regexp) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 检查是否为白名单IP
		if skip, _ := c.Get("isWhiteListIP"); skip == true {
			c.Next()
			return
		}

		// 检查URL查询参数和表单数据中的SQL注入
		query := c.Request.URL.Query()
		for _, values := range query {
			for _, value := range values {
				for _, pattern := range rules {
					if pattern.MatchString(value) {
						c.Set("BlockedBy", "Sql注入")
						c.Set("BlockReason", "检查到存在Sql注入")
						api_handler.ForbiddenHandler(c, "检测到sql注入攻击，禁止访问！")
						c.Abort()
						return
					}
				}
			}
		}

		postForm, _ := c.GetRawData()
		for _, pattern := range rules {
			if pattern.MatchString(string(postForm)) {
				c.Set("BlockedBy", "Sql注入")
				c.Set("BlockReason", "检查到存在Sql注入")
				api_handler.ForbiddenHandler(c, "检测到sql注入攻击，禁止访问！")
				c.Abort()
				return
			}
		}

		c.Next()
	}
}
