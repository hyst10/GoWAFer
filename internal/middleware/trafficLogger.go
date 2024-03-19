package middleware

import (
	"GoWAFer/internal/model"
	"GoWAFer/internal/repository"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"time"
)

// TrafficLogger 用于记录流量日志的中间件
func TrafficLogger(r *repository.LogRepository) gin.HandlerFunc {
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

		// 记录请求日志
		currentLog := model.Log{
			ClientIP:    requestIP,
			Method:      requestMethod,
			Url:         requestPath,
			UserAgent:   c.GetHeader("User-Agent"),
			Referer:     c.GetHeader("Referer"),
			ContentType: c.ContentType(),
			Latency:     latency,
			Status:      responseStatus,
		}

		// 检查是否被WAF中间件拦截,被拦截则保留更多请求细节。
		blockedBy, exists := c.Get("BlockedBy")
		if exists {
			// 获取拦截原因
			blockReason, _ := c.Get("BlockReason")

			currentLog.UserAgent = c.GetHeader("User-Agent")
			currentLog.Referer = c.GetHeader("Referer")
			currentLog.ContentType = c.ContentType()
			currentLog.BlockBy = blockedBy.(string)
			currentLog.BlockReason = blockReason.(string)
		}
		// 保存日志到数据库
		err := r.Create(&currentLog)
		if err != nil {
			log.Printf("日志保存失败：%v", err)
		}

	}
}
