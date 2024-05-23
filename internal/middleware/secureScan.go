package middleware

import (
	"GoWAFer/constants"
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net/url"
)

// SecureRequestMiddleware 参数安全检测中间件
func SecureRequestMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 检查是否为白名单IP
		if skip, _ := c.Get("isWhiteIP"); skip == true {
			c.Next()
			return
		}
		// 检查是否为白名单路径
		if skip, _ := c.Get("isWhitePath"); skip == true {
			c.Next()
			return
		}

		rawQuery := c.Request.URL.RawQuery
		decodedQuery, err := url.QueryUnescape(rawQuery)
		if err != nil {
			log.Printf("url解码失败：%v", err)
		}
		fmt.Printf("参数：%s\n", decodedQuery)
		for pattern := range constants.SqlInjectRules {
			if pattern.MatchString(decodedQuery) {
				c.Set("BlockedBy", "sql注入拦截中间件")
				c.Set("BlockReason", "路径参数中存在sql注入")
				c.HTML(200, "block.html", gin.H{"Reason": "路径参数中存在sql注入"})
				c.Abort()
				return
			}
		}
		for pattern := range constants.XssDetectRules {
			if pattern.MatchString(decodedQuery) {
				c.Set("BlockedBy", "xss攻击拦截中间件")
				c.Set("BlockReason", "路径参数中存在xss攻击")
				c.HTML(200, "block.html", gin.H{"Reason": "路径参数中存在xss攻击"})
				c.Abort()
				return
			}
		}

		// 检查表单数据
		contentType := c.GetHeader("Content-Type")
		if contentType == "application/x-www-form-urlencoded" || contentType == "multipart/form-data" {
			if err := c.Request.ParseForm(); err == nil {
				for _, values := range c.Request.PostForm {
					for _, value := range values {
						for pattern := range constants.XssDetectRules {
							if pattern.MatchString(value) {
								c.Set("BlockedBy", "xss攻击拦截中间件")
								c.Set("BlockReason", "表单数据中存在xss攻击")
								c.HTML(200, "block.html", gin.H{"Reason": "表单数据中存在xss攻击"})
								c.Abort()
								return
							}
						}
					}
				}
			}
		}

		body, _ := c.GetRawData()
		for pattern := range constants.SqlInjectRules {
			if pattern.MatchString(decodedQuery) {
				c.Set("BlockedBy", "sql注入中间件")
				c.Set("BlockReason", "请求体中存在ql注入")
				c.HTML(200, "block.html", gin.H{"Reason": "请求体中存在ql注入"})
				c.Abort()
				return
			}
		}

		c.Request.Body = io.NopCloser(bytes.NewBuffer(body))
		c.Next()
	}
}
