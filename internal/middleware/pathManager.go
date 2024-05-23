package middleware

import (
	"GoWAFer/constants"
	"GoWAFer/internal/repository"
	"github.com/gin-gonic/gin"
)

// PathManager 路径管理中间件
func PathManager(r *repository.RoutingManageRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 检查是否为白名单IP
		if skip, _ := c.Get("isWhiteIP"); skip == true {
			c.Next()
			return
		}
		urlPath := c.Request.URL.Path
		exist := r.IsExist(urlPath)
		if exist == constants.BlackPathKey {
			c.Set("BlockedBy", "路径管理中间件")
			c.Set("BlockReason", "黑名单路径禁止访问")
			c.HTML(200, "block.html", gin.H{"Reason": "黑名单路径禁止访问"})
			c.Abort()
			return
		} else if exist == constants.WhitePathKey {
			c.Set("isWhitePath", true)
		}
		c.Next()
		return
	}
}
