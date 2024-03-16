package middleware

import (
	"GoWAFer/internal/model"
	"GoWAFer/internal/repository"
	"fmt"
	"github.com/gin-gonic/gin"
	"time"
)

// TrafficLogger 记录访问日志的中间件
func TrafficLogger(r *repository.LogRepository, r2 *repository.BlockLogRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 请求处理前
		startTime := time.Now()

		// 请求信息
		requestPath := c.Request.URL.Path
		rawQuery := c.Request.URL.RawQuery
		if rawQuery != "" {
			requestPath = fmt.Sprintf("%s?%s", requestPath, rawQuery)
		}
		requestMethod := c.Request.Method
		requestIP := c.ClientIP()

		c.Next()

		responseStatus := c.Writer.Status()
		endTime := time.Now()
		latency := endTime.Sub(startTime)

		// 检查是否被WAF中间件拦截
		blockedBy, exists := c.Get("BlockedBy")
		if exists {
			blockReason, _ := c.Get("BlockReason")
			log := model.Log{
				Path:    requestPath,
				Method:  requestMethod,
				IP:      requestIP,
				Latency: latency,
				Status:  responseStatus,
			}
			r.Create(&log)

			blockLog := model.BlockLog{
				ClientIP:    requestIP,
				Method:      requestMethod,
				Url:         requestPath,
				UserAgent:   c.GetHeader("User-Agent"),
				Referer:     c.GetHeader("Referer"),
				ContentType: c.ContentType(),
				Latency:     latency,
				BlockBy:     blockedBy.(string),
				BlockReason: blockReason.(string),
			}
			r2.Create(&blockLog)
		}

	}
}
