package middleware

import (
	"GoWAFer/internal/repository"
	"GoWAFer/pkg/utils/api_helper"
	"github.com/gin-gonic/gin"
)

// IPManager IP管理中间件，用于检验IP黑白名单的中间件
func IPManager(r *repository.IPManageRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取客户端IP
		clientIP := c.ClientIP()
		// 检查是否为黑名单IP
		isExist, _ := r.IsExist(clientIP, 1)
		if isExist {
			c.Set("BlockedBy", "IP管理中间件")
			c.Set("BlockReason", "客户端IP为黑名单IP")
			api_helper.ForbiddenHandler(c, "该IP已被添加至黑名单，禁止访问！")
			c.Abort()
			return
		}
		// 检查是否为白名单IP
		isExist, _ = r.IsExist(clientIP, 2)
		if isExist {
			c.Set("isWhiteListIP", true)
			c.Next()
			return
		}
		c.Next()
		return
	}
}
