package middleware

import (
	"GoWAFer/internal/repository"
	"GoWAFer/pkg/utils/api_helper"
	"github.com/gin-gonic/gin"
)

// RouteGuardMiddleware 路由守卫中间件
func RouteGuardMiddleware(r *repository.RoutingManageRepository) gin.HandlerFunc {
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
		// 检查路由是否为白名单路由
		exist, _ := r.IsExist(path, method, 2)
		if exist {
			c.Set("isWhiteRoute", true)
			c.Next()
			return
		}
		// 检查路由是否为黑名单路由
		exist, _ = r.IsExist(path, method, 1)
		if exist {
			c.Set("BlockedBy", "路由守卫中间件")
			c.Set("BlockReason", "路由为黑名单路由")
			api_helper.ForbiddenHandler(c, "此路由不可访问！")
			c.Abort()
			return
		}
		c.Next()
		return
	}
}
