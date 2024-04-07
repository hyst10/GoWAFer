package middleware

import (
	"GoWAFer/internal/types"
	"GoWAFer/pkg/utils/api_helper"
	"bytes"
	"github.com/gin-gonic/gin"
	"io"
)

// SqlInjectMiddleware sql注入检测中间件
func SqlInjectMiddleware() gin.HandlerFunc {
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

		// 检查URL查询参数和body中是否存在SQL注入
		query := c.Request.URL.Query()
		for _, values := range query {
			for _, value := range values {
				for pattern := range types.SqlInjectRules {
					if pattern.MatchString(value) {
						c.Set("BlockedBy", "sql注入防护中间件")
						c.Set("BlockReason", "查询参数中检测到sql注入")
						api_helper.ForbiddenHandler(c, "检测到sql注入攻击，禁止访问！")
						c.Abort()
						return
					}
				}
			}
		}

		// c.GetRawData()会读取并消耗掉http.Request的Body，这意味着Body流被读取后，如果不进行特殊处理，就无法再次读取。这在后续的处理中可能会导致问题
		body, _ := c.GetRawData()
		for pattern := range types.SqlInjectRules {
			if pattern.MatchString(string(body)) {
				c.Set("BlockedBy", "sql注入防护中间件")
				c.Set("BlockReason", "请求体body中检测到sql注入")
				api_helper.ForbiddenHandler(c, "检测到sql注入攻击，禁止访问！")
				c.Abort()
				return
			}
		}

		// 将原始请求体重新写回
		c.Request.Body = io.NopCloser(bytes.NewBuffer(body))
		c.Next()
	}
}
