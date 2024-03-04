package main

import (
	"GoWAFer/api"
	_ "GoWAFer/docs"
	"GoWAFer/pkg/config"
	"GoWAFer/pkg/database"
	"fmt"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title GoWAFer
// @description Golang编写的一款基于反向代理模式的web防火墙应用 By supercat0867
// @version v0.1
func main() {
	// 读取配置文件
	conf := config.ReadConfig()

	// 创建数据库连接
	db, err := database.InitDB()
	if err != nil {
		panic(fmt.Sprintf("数据库连接失败：%v", err))
	}

	r := gin.Default()

	// swagger 文档
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	// 注册API路由
	api.RegisterAllHandlers(r, db, conf)
	// 注册页面路由

	// 启动waf服务
	r.Run(fmt.Sprintf(":%d", conf.Server.Port))

}
