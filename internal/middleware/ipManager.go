package middleware

import (
	"GoWAFer/constants"
	"GoWAFer/internal/repository"
	"github.com/gin-gonic/gin"
)

// IPManager IP管理中间件
func IPManager(r *repository.IPManageRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		clientIP := c.ClientIP()
		isExist := r.IsExist(clientIP)
		if isExist == constants.BlackIPKey {
			c.Set("BlockedBy", "IP管理中间件")
			c.Set("BlockReason", "黑名单IP禁止访问")
			c.HTML(200, "block.html", gin.H{"Reason": "黑名单IP禁止访问"})
			c.Abort()
			return
		} else if isExist == constants.WhiteIPKey {
			c.Set("isWhiteIP", true)
		}
		c.Next()
		return
	}
}
