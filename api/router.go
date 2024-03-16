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
	"regexp"
)

type Databases struct {
	userRepository      *repository.UserRepository
	logRepository       *repository.LogRepository
	ipRepository        *repository.IPRepository
	blockRepository     *repository.BlockLogRepository
	sqlInjectRepository *repository.SqlInjectRepository
}

func NewDatabases(db *gorm.DB) *Databases {
	return &Databases{
		userRepository:      repository.NewUserRepository(db),
		logRepository:       repository.NewLogRepository(db),
		ipRepository:        repository.NewIPRepository(db),
		blockRepository:     repository.NewBlockLogRepository(db),
		sqlInjectRepository: repository.NewSqlInjectRepository(db),
	}
}

func RegisterAllHandlers(r *gin.Engine, db *gorm.DB, conf *config.Config) {
	dbs := NewDatabases(db)

	apiGroup := r.Group("/waf/api/v1")

	RegisterUserHandler(apiGroup, dbs, conf)
	RegisterIPHandler(apiGroup, dbs, conf)
	RegisterLogHandler(apiGroup, dbs, conf)
	RegisterConfigHandler(apiGroup, dbs, conf)
	RegisterBlockLogHandler(apiGroup, dbs, conf)
	RegisterSqlInjectHandler(apiGroup, dbs, conf)

	// 添加日志中间件
	r.Use(middleware.TrafficLogger(dbs.logRepository, dbs.blockRepository))
	// 添加IP管理中间件
	r.Use(middleware.IPManager(dbs.ipRepository))
	// 添加限速器中间件
	r.Use(middleware.RateLimitMiddleware(conf, dbs.ipRepository))
	// 加载sql注入规则
	sqlInjectRules := dbs.sqlInjectRepository.FindAll()
	var sqlInject []*regexp.Regexp
	for _, rule := range sqlInjectRules {
		compile := regexp.MustCompile(rule.Rule)
		sqlInject = append(sqlInject, compile)
	}
	// 添加sql注入检测中间件
	r.Use(middleware.SqlInjectMiddleware(sqlInject))

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

func RegisterBlockLogHandler(r *gin.RouterGroup, dbs *Databases, conf *config.Config) {
	blockLogService := service.NewBlockLogService(dbs.blockRepository)
	blockLogController := controller.NewBlockLogController(blockLogService)
	blockLogGroup := r.Group("/blockLog")
	blockLogGroup.Use(middleware.WafAPIAuthMiddleware(conf.Secret.JwtSecretKey, dbs.userRepository))
	blockLogGroup.GET("", blockLogController.FindPaginatedBlockLog)
}

func RegisterSqlInjectHandler(r *gin.RouterGroup, dbs *Databases, conf *config.Config) {
	sqlInjectService := service.NewSqlInjectService(dbs.sqlInjectRepository)
	sqlInjectController := controller.NewSqlInjectController(sqlInjectService)
	sqlInjectGroup := r.Group("/sqlInject")
	sqlInjectGroup.Use(middleware.WafAPIAuthMiddleware(conf.Secret.JwtSecretKey, dbs.userRepository))
	sqlInjectGroup.GET("", sqlInjectController.FindPaginatedSqlInject)
	sqlInjectGroup.POST("", sqlInjectController.CreateRule)
	sqlInjectGroup.PATCH(":id", sqlInjectController.UpdateRule)
	sqlInjectGroup.DELETE(":id", sqlInjectController.DeleteRule)
}
