package model

import (
	"gorm.io/gorm"
	"time"
)

// Admin 管理员模型
type Admin struct {
	gorm.Model
	Username      string    // 用户名
	Password      string    // 密码
	Nickname      string    // 昵称
	LastLoginDate time.Time // 上次登录日期
	LastLoginIP   string    // 上次登录IP
	RefreshToken  string    // refreshToken
}

// Log 日志模型
type Log struct {
	ID          uint          `gorm:"primarykey" json:"id"`
	ClientIP    string        `json:"clientIP"`  // 客户端IP
	Method      string        `json:"method"`    // http方法
	Url         string        `json:"url"`       // 完整的请求url
	UserAgent   string        `json:"userAgent"` // ua头
	Referer     string        `json:"referer"`   // 来源页面
	ContentType string        `json:"contentType"`
	BlockBy     string        `json:"blockBy"`
	BlockReason string        `json:"blockReason"`
	Latency     time.Duration `json:"latency"`
	Status      int           `json:"status"`
	CreatedAt   time.Time     `json:"createdAt"`
}
