package middleware

import (
	"GoWAFer/pkg/utils/api_helper"
	"GoWAFer/pkg/utils/jwt_handler"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// CsrfTokenMiddleware csrfToken中间件
func CsrfTokenMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		session := sessions.Default(c)

		if c.Request.Method != "GET" {
			// 验证CSRF令牌
			//从请求头中获取csrf-token
			csrfTokenFromRequest := c.GetHeader("X-CSRF-Token")
			// 从用户会话中获取CSRF令牌
			csrfTokenFromSession := session.Get("csrfToken")
			if csrfTokenFromRequest != csrfTokenFromSession {
				// CSRF令牌验证失败
				c.Set("BlockedBy", "CSRF-Token中间件")
				c.Set("BlockReason", "CSRFToken验证失败")
				api_helper.ForbiddenHandler(c, "CSRFToken验证失败！")
				c.Abort()
				return
			}
		}

		// 生成新的CSRF令牌
		csrfToken, _ := jwt_handler.GenerateRandomKey(32)
		// 存储CSRF令牌在用户会话中
		session.Set("csrfToken", csrfToken)
		session.Save()
		// 将CSRF令牌添加到响应头，以便客户端JavaScript代码可以访问
		c.Header("X-CSRF-Token", csrfToken)
		c.Next()
		return
	}
}
