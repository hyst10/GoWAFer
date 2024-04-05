package middleware

import (
	"GoWAFer/internal/types"
	"GoWAFer/pkg/utils/api_helper"
	"errors"
	"github.com/gin-gonic/gin"
	"log"
)

// ErrorHandlingMiddleware 错误处理层中间件
func ErrorHandlingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {

			err := c.Errors[0].Err
			var httpStatusCode = 200 // http状态码
			var errMsg = err.Error() // 错误信息
			var errCode int          // 业务码

			switch {
			case errors.Is(err, types.ErrInvalidBody):
				errCode = 40001
			case errors.Is(err, types.ErrRoutingNotMatch):
				errCode = 40002
			case errors.Is(err, types.ErrIPNotMatch):
				errCode = 40005
			default:
				log.Printf("服务器异常错误：%v\n", err)
				httpStatusCode = 500
				errCode = 50000
				errMsg = "服务器开小差啦，请稍后再试！"
			}

			c.JSON(httpStatusCode, api_helper.Response{Status: errCode, Msg: errMsg})
			c.Abort()
		}
	}
}
