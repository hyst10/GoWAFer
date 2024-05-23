package middleware

import (
	"GoWAFer/internal/model"
	"GoWAFer/internal/repository"
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"time"
)

// TrafficLogger 用于记录流量日志的中间件
func TrafficLogger(r *repository.LogRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 计算处理时间
		startTime := time.Now()
		// 保存原始请求体以备恢复
		originalBody := c.Request.Body
		// 读取并转换请求体为字符串
		bodyBytes, err := ioutil.ReadAll(originalBody)
		if err != nil {
			c.Error(fmt.Errorf("请求体读取失败：%v", err))
			c.Abort()
			return
		}
		bodyStr := string(bodyBytes)
		// 恢复请求体
		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

		c.Next()

		newLog := model.Log{
			IP:        c.ClientIP(),
			Path:      c.Request.URL.Path,
			Query:     c.Request.URL.RawQuery,
			Body:      bodyStr,
			UserAgent: c.Request.UserAgent(),
			Referer:   c.Request.Referer(),
			Method:    c.Request.Method,
			Status:    c.Writer.Status(),
			Latency:   time.Since(startTime).Milliseconds(),
		}
		blockedBy, exists := c.Get("BlockedBy")
		if exists {
			blockReason, _ := c.Get("BlockReason")
			newLog.BlockBy = blockedBy.(string)
			newLog.BlockReason = blockReason.(string)
		}
		if err := r.Create(&newLog); err != nil {
			c.Error(fmt.Errorf("日志保存失败：%v", err))
			c.Abort()
			return
		}
	}
}
