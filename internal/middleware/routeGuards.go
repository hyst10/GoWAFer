package middleware

import (
	"GoWAFer/internal/repository"
	"GoWAFer/pkg/utils/api_handler"
	"github.com/gin-gonic/gin"
)

// RouteGuardMiddleware 路由守卫中间件
func RouteGuardMiddleware(r *repository.RoutingRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 检查是否为白名单IP
		if skip, _ := c.Get("isWhiteListIP"); skip == true {
			c.Next()
			return
		}

		// 获取当前路由
		path := c.Request.URL.Path
		// 获取Method
		method := c.Request.Method

		// 检查是否存在路由表中
		current, err := r.IsExist(path, method)
		if err != nil {
			c.Next()
			return
		}

		if current.Type == 1 {
			c.Set("isWhiteRoute", true)
			c.Next()
			return
		}
		// 黑名单路由
		c.Set("BlockedBy", "路由守卫中间件")
		c.Set("BlockReason", "路由为黑名单路由")
		api_handler.ForbiddenHandler(c, "此路由不可访问！")
		c.Abort()
		return
	}
}
