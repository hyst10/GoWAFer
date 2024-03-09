package model

import (
	"gorm.io/gorm"
	"time"
)

// User 管理员模型
type User struct {
	gorm.Model
	Username      string    `gorm:"comment:用户名"`
	Password      string    `gorm:"comment:密码"`
	Nickname      string    `gorm:"comment:昵称"`
	LastLoginDate time.Time `gorm:"comment:上次登录日期"`
	LastLoginIP   string    `gorm:"comment:上次登录IP"`
	RefreshToken  string    `gorm:"comment:refreshToken"`
}

// IP IP管理模型
type IP struct {
	ID           uint      `gorm:"primarykey" json:"id"` // 自增主键
	IPAddress    string    `json:"ip_address"`           // IP地址
	Type         int       `json:"type"`                 // 类型 1:黑名单 2:白名单
	CreatedAt    time.Time `json:"created_at"`           // 添加日期
	ExpirationAt time.Time `json:"expiration_at"`        // 过期时间
}
