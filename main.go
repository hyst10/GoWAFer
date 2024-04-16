package main

import (
	"GoWAFer/api"
	_ "GoWAFer/docs"
	"GoWAFer/migrate"
	"GoWAFer/pkg/config"
	"GoWAFer/pkg/database"
	"GoWAFer/pkg/utils/graceful"
	"GoWAFer/web"
	"context"
	"errors"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// @title GoWAFer
// @description Golang编写的一款基于反向代理模式的web防火墙应用 By supercat0867
// @version v0.1
func main() {
	graceful.Welcome()

	// 读取配置文章
	conf := config.ReadConfig()
	log.Println("配置文件加载成功")

	// 创建数据库连接
	db, err := database.InitDB()
	if err != nil {
		panic(fmt.Sprintf("数据库连接失败：%v", err))
	}
	log.Println("数据库连接成功")
	rdb := database.InitRedis()

	// 迁移数据库
	migrate.AutoMigrateAndInsertData(db)
	log.Println("数据库迁移成功")

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
	api.RegisterAllHandlers(r, db, rdb, conf)

	// 创建一个监听系统终止信号的通道
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, syscall.SIGINT, syscall.SIGTERM)

	// 启动waf服务
	srv := &http.Server{Addr: ":8080", Handler: r}
	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			fmt.Printf("listen: %s\n", err)
		}
	}()

	// 阻塞，直到接收到系统终止信号
	<-stopChan
	fmt.Println("Shutting down server...")

	// 执行清理操作
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		fmt.Printf("Server Shutdown: %v", err)
	}

	fmt.Println("Server gracefully stopped")
}
