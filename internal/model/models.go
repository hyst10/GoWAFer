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
