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
	ID          uint      `gorm:"primarykey;autoIncrement" json:"id"`
	IP          string    `gorm:"size:200;not null" json:"ip"`    // 客户端IP
	Method      string    `gorm:"size:50;not null" json:"method"` // 请求方法
	Path        string    `gorm:"size:255;not null" json:"path"`  // 请求路径
	Query       string    `gorm:"size:255"  json:"query"`         // 请求参数
	Body        string    `gorm:"type:text" json:"body"`          // 请求体
	UserAgent   string    `gorm:"size:255" json:"userAgent"`      // 请求头
	Referer     string    `gorm:"size:255" json:"referer"`        // 请求来源
	Status      int       `gorm:"not null" json:"status"`         // 状态码
	Latency     int64     `gorm:"not null" json:"latency"`        // 耗时(毫秒)
	BlockBy     string    `gorm:"size:50" json:"blockBy"`         // 拦截中间件
	BlockReason string    `gorm:"size:255" json:"blockReason"`    // 拦截原因
	CreatedAt   time.Time `gorm:"not null" json:"createdAt"`      // 请求时间
}
