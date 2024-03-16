package main

import (
	"GoWAFer/internal/model"
	"GoWAFer/pkg/database"
	"GoWAFer/pkg/hash_handler"
	"fmt"
	"gorm.io/gorm"
)

func migrate(db *gorm.DB) {
	err := db.AutoMigrate(&model.User{})
	if err != nil {
		panic(fmt.Sprintf("管理员表创建失败：%v", err))
	}

	err = db.AutoMigrate(&model.IP{})
	if err != nil {
		panic(fmt.Sprintf("IP管理表创建失败：%v", err))
	}

	err = db.AutoMigrate(&model.Log{})
	if err != nil {
		panic(fmt.Sprintf("流量日志表创建失败：%v", err))
	}

	err = db.AutoMigrate(&model.BlockLog{})
	if err != nil {
		panic(fmt.Sprintf("拦截日志表创建失败：%v", err))
	}

	err = db.AutoMigrate(&model.SqlInjectionRules{})
	if err != nil {
		panic(fmt.Sprintf("sql注入规则表创建失败：%v", err))
	}
}

func insert(db *gorm.DB) {
	// 哈希加盐加密
	defaultPassword, _ := hash_handler.EncryptPassword("123456")
	defaultUser := model.User{Model: gorm.Model{ID: 1}, Username: "admin", Password: defaultPassword, Nickname: "admin"}
	err := db.FirstOrCreate(&defaultUser, model.User{Model: gorm.Model{ID: 1}}).Error
	if err != nil {
		panic(fmt.Sprintf("默认管理员创建失败：%v", err))
	}

	sqlDefault := []model.SqlInjectionRules{
		{Rule: "(?i)(union)(.*)(select)"},
		{Rule: "(?i)select(.*)from"},
		{Rule: "(?i)insert into"},
		{Rule: "(?i)delete from"},
		{Rule: "(?i)drop table"},
		{Rule: "(?i)update(.*)set"},
		{Rule: "--"},
		{Rule: "(\\b|\\')(OR|or|oR|Or)('|\\b)\\s*('\\d+'|'\\d+'--\\s*|'\\d+'(\\s+)(--)?|\\d+)(\\s+)(=|like)(\\s+)(\\b|\\')\\d+('|\\b)"},
		{Rule: "/\\*.*\\*/"},
		{Rule: ";"},
	}
	for _, sql := range sqlDefault {
		db.FirstOrCreate(&sql, model.SqlInjectionRules{Rule: sql.Rule})
	}
}

func main() {
	// 创建数据库连接
	db, err := database.InitDB()
	if err != nil {
		panic(fmt.Sprintf("数据库连接失败：%v", err))
	}

	migrate(db)
	insert(db)
}
