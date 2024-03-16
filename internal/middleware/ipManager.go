package middleware

import (
	"GoWAFer/internal/repository"
	"GoWAFer/pkg/utils/api_handler"
	"github.com/gin-gonic/gin"
)

// IPManager IP管理中间件，用于检验IP黑白名单的中间件
func IPManager(r *repository.IPRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取客户端IP
		clientIP := c.ClientIP()
		// 检查是否在IP名单中
		current, err := r.IsIPExist(clientIP)
		if err != nil {
			c.Next()
			return
		}
		// 存在名单中，检查IP为黑名单还是白名单
		if current.Type == 1 {
			// 白名单
			// 为白名单IP上下文设置一个标识，用于绕过后续的中间件
			c.Set("isWhiteListIP", true)
			c.Next()
			return
		}
		// 黑名单
		c.Set("BlockedBy", "IP管理")
		c.Set("BlockReason", "客户端IP为黑名单IP")
		api_handler.ForbiddenHandler(c, "该IP已被添加至黑名单，禁止访问！")
		c.Abort()
		return
	}
}
