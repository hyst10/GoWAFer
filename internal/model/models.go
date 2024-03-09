package model

import (
	"gorm.io/gorm"
	"time"
)

// User 管理员模型
type User struct {
	gorm.Model
	Username      string    // 用户名
	Password      string    // 密码
	Nickname      string    // 昵称
	LastLoginDate time.Time // 上次登录日期
	LastLoginIP   string    // 上次登录IP
	RefreshToken  string    // refreshToken
}

// IP IP模型
type IP struct {
	ID           uint      `gorm:"primarykey" json:"id"` // 自增主键
	IPAddress    string    `json:"ip_address"`           // IP地址
	Type         int       `json:"type"`                 // 类型 1:黑名单 2:白名单
	CreatedAt    time.Time `json:"created_at"`           // 添加日期
	ExpirationAt time.Time `json:"expiration_at"`        // 过期时间
}

// Log 日志模型
type Log struct {
	ID        uint   `gorm:"primarykey" json:"id"`
	IP        string // IP
	Path      string // Path
	Method    string // http方法
	Status    int
	Latency   time.Duration
	CreatedAt time.Time
}
