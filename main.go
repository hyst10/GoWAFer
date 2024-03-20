package main

import (
	"GoWAFer/api"
	_ "GoWAFer/docs"
	"GoWAFer/internal/repository"
	"GoWAFer/migrate"
	"GoWAFer/pkg/config"
	"GoWAFer/pkg/database"
	"GoWAFer/pkg/utils/graceful"
	"GoWAFer/web"
	"context"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
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

	// 读取配置文件
	conf := config.ReadConfig()
	log.Println("配置文件加载成功")

	// 创建数据库连接
	db, err := database.InitDB()
	if err != nil {
		panic(fmt.Sprintf("数据库连接失败：%v", err))
	}
	log.Println("数据库连接成功")

	// 创建一个协程用于管理过期的黑白名单IP
	go func() {
		repo := repository.NewIPRepository(db)
		// 创建一个每秒触发一次的定时器
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()

		// 无限循环，等待定时器的信号
		for {
			select {
			case <-ticker.C:
				// 每当定时器触发，调用DeleteExpired函数
				err := repo.DeleteExpired()
				if err != nil {
					// 处理错误
					fmt.Println("Error deleting expired records:", err)
				}
			}
		}
	}()

	// 迁移数据
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
	//r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// 注册页面路由
	web.RegisterWebHandler(r, db, conf)

	// 注册API路由、注册反向代理路由
	api.RegisterAllHandlers(r, db, conf)

	// 创建一个监听系统终止信号的通道
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, syscall.SIGINT, syscall.SIGTERM)

	// 启动waf服务
	srv := &http.Server{Addr: ":8080", Handler: r}
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
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
