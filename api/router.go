package api

import (
	"GoWAFer/internal/controller"
	"GoWAFer/internal/middleware"
	"GoWAFer/internal/repository"
	"GoWAFer/internal/service"
	"GoWAFer/pkg/config"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http/httputil"
	"net/url"
)

type Databases struct {
	userRepository *repository.UserRepository
	logRepository  *repository.LogRepository
	ipRepository   *repository.IPRepository
}

func NewDatabases(db *gorm.DB) *Databases {
	return &Databases{
		userRepository: repository.NewUserRepository(db),
		logRepository:  repository.NewLogRepository(db),
		ipRepository:   repository.NewIPRepository(db),
	}
}

func RegisterAllHandlers(r *gin.Engine, db *gorm.DB, conf *config.Config) {
	dbs := NewDatabases(db)

	apiGroup := r.Group("/waf/api/v1")

	RegisterUserHandler(apiGroup, dbs, conf)
	RegisterIPHandler(apiGroup, dbs, conf)
	RegisterLogHandler(apiGroup, dbs, conf)
	RegisterConfigHandler(apiGroup, dbs, conf)

	// 添加日志中间件
	r.Use(middleware.TrafficLogger(dbs.logRepository))
	// 添加IP管理中间件
	r.Use(middleware.IPManager(dbs.ipRepository))
	// 添加限速器中间件
	r.Use(middleware.RateLimitMiddleware(conf, dbs.ipRepository))

	// 设置反向代理
	target, _ := url.Parse(conf.Server.TargetAddress)
	proxy := httputil.NewSingleHostReverseProxy(target)
	r.NoRoute(func(c *gin.Context) {
		proxy.ServeHTTP(c.Writer, c.Request)
	})
}

func RegisterUserHandler(r *gin.RouterGroup, dbs *Databases, conf *config.Config) {
	userService := service.NewUserService(dbs.userRepository)
	userController := controller.NewUserController(userService, conf)
	authGroup := r.Group("/auth")
	authGroup.POST("/dologin", userController.DoLogin)
	authGroup.GET("/getCaptcha", userController.GetCaptcha)
	// 中间件
}

func RegisterIPHandler(r *gin.RouterGroup, dbs *Databases, conf *config.Config) {
	ipService := service.NewIPListService(dbs.ipRepository)
	ipController := controller.NewIPController(ipService)
	ipGroup := r.Group("/ip")
	ipGroup.Use(middleware.WafAPIAuthMiddleware(conf.Secret.JwtSecretKey, dbs.userRepository))
	ipGroup.POST("", ipController.CreateIP)
	ipGroup.GET("", ipController.FindPaginatedIP)
	ipGroup.PATCH(":id", ipController.UpdateIP)
	ipGroup.DELETE(":id", ipController.DeleteIP)
}

func RegisterLogHandler(r *gin.RouterGroup, dbs *Databases, conf *config.Config) {
	logService := service.NewLogService(dbs.logRepository)
	logController := controller.NewLogController(logService)
	logGroup := r.Group("/log")
	logGroup.Use(middleware.WafAPIAuthMiddleware(conf.Secret.JwtSecretKey, dbs.userRepository))
	logGroup.GET("", logController.FindLogs)
}

func RegisterConfigHandler(r *gin.RouterGroup, dbs *Databases, conf *config.Config) {
	configController := controller.NewConfigController(conf)
	configGroup := r.Group("/config")
	configGroup.Use(middleware.WafAPIAuthMiddleware(conf.Secret.JwtSecretKey, dbs.userRepository))
	configGroup.GET("", configController.GetConfig)
	configGroup.PUT("", configController.UpdateConfig)
}
