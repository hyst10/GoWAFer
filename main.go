package main

import (
	"GoWAFer/api"
	_ "GoWAFer/docs"
	"GoWAFer/migrate"
	"GoWAFer/pkg/config"
	"GoWAFer/pkg/database"
	"GoWAFer/pkg/utils/graceful"
	"GoWAFer/web"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"log"
	"net/http"
	"time"
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

	// 迁移数据
	migrate.AutoMigrateAndInsertData(db)

	r := gin.Default()

	store := cookie.NewStore([]byte(conf.Secret.SessionSecretKey))
	store.Options(sessions.Options{
		MaxAge:   60 * 60 * 24 * 30,
		Path:     "/",
		HttpOnly: true,
		//Secure:   true,
	})
	r.Use(sessions.Sessions("waf-session", store))

	// 加载静态资源
	r.Static("/static", "./static")
	// 加载模版文件
	r.LoadHTMLGlob("templates/**")

	// swagger 文档
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	// 注册页面路由
	web.RegisterWebHandler(r, db, conf)
	// 注册API路由、注册反向代理路由
	api.RegisterAllHandlers(r, db, conf)

	// 启动waf服务
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", conf.Server.WafPort),
		Handler: r,
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Printf("服务异常：%s\n", err)
		}
	}()

	graceful.Welcome()
	graceful.ShutdownGin(srv, time.Second*3)
}
