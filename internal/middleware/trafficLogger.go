package middleware

import (
	"GoWAFer/internal/model"
	"GoWAFer/internal/repository"
	"github.com/gin-gonic/gin"
	"time"
)

// TrafficLogger 记录访问日志的中间件
func TrafficLogger(r *repository.LogRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 请求处理前
		startTime := time.Now()

		// 请求信息
		requestPath := c.Request.URL.Path
		requestMethod := c.Request.Method
		requestIP := c.ClientIP()

		c.Next()

		responseStatus := c.Writer.Status()
		endTime := time.Now()
		latency := endTime.Sub(startTime)

		log := model.Log{
			Path:    requestPath,
			Method:  requestMethod,
			IP:      requestIP,
			Latency: latency,
			Status:  responseStatus,
		}
		r.Create(&log)
	}
}
