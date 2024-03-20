package api

import (
	"GoWAFer/internal/controller"
	"GoWAFer/internal/middleware"
	"GoWAFer/internal/repository"
	"GoWAFer/internal/service"
	"GoWAFer/pkg/config"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log"
	"net/http/httputil"
	"net/url"
	"regexp"
)

type Databases struct {
	adminRepository     *repository.AdminRepository
	logRepository       *repository.LogRepository
	ipRepository        *repository.IPRepository
	routingRepository   *repository.RoutingRepository
	sqlInjectRepository *repository.SqlInjectRepository
	xssDetectRepository *repository.XssDetectRepository
}

func NewDatabases(db *gorm.DB) *Databases {
	return &Databases{
		adminRepository:     repository.NewAdminRepository(db),
		logRepository:       repository.NewLogRepository(db),
		ipRepository:        repository.NewIPRepository(db),
		routingRepository:   repository.NewRoutingRepository(db),
		sqlInjectRepository: repository.NewSqlInjectRepository(db),
		xssDetectRepository: repository.NewXssDetectRepository(db),
	}
}

func RegisterAllHandlers(r *gin.Engine, db *gorm.DB, conf *config.Config) {
	dbs := NewDatabases(db)

	apiGroup := r.Group("/waf/api/v1")

	RegisterUserHandler(apiGroup, dbs, conf)
	RegisterIPHandler(apiGroup, dbs, conf)
	RegisterRoutingHandler(apiGroup, dbs, conf)
	RegisterLogHandler(apiGroup, dbs, conf)
	RegisterConfigHandler(apiGroup, dbs, conf)
	RegisterSqlInjectHandler(apiGroup, dbs, conf)
	RegisterXssDetectHandler(apiGroup, dbs, conf)

	// 添加日志中间件
	r.Use(middleware.TrafficLogger(dbs.logRepository))
	log.Println("日志中间件加载成功")
	// 添加IP管理中间件
	r.Use(middleware.IPManager(dbs.ipRepository))
	log.Println("IP管理中间件加载成功")
	// 添加路由守卫中间件
	r.Use(middleware.RouteGuardMiddleware(dbs.routingRepository))
	log.Println("路由守卫中间件加载成功")
	// 添加CSRFToken中间件
	r.Use(middleware.CsrfTokenMiddleware())
	log.Println("CSRFToken中间件加载成功")
	// 添加限速器中间件
	r.Use(middleware.RateLimitMiddleware(conf, dbs.ipRepository))
	log.Println("CC攻击防护中间件加载成功")
	// 加载sql注入防护规则
	sqlInjectRules := dbs.sqlInjectRepository.FindAll()
	var sqlInject []*regexp.Regexp
	for _, rule := range sqlInjectRules {
		compile := regexp.MustCompile(rule.Rule)
		sqlInject = append(sqlInject, compile)
	}
	log.Println("sql注入防护规则加载成功")
	// 添加sql注入检测中间件
	r.Use(middleware.SqlInjectMiddleware(sqlInject))
	log.Println("sql注入防护中间件加载成功")
	// 加载xss攻击防护规则
	xssDetectRules := dbs.xssDetectRepository.FindAll()
	var xssDetect []*regexp.Regexp
	for _, rule := range xssDetectRules {
		compile := regexp.MustCompile(rule.Rule)
		xssDetect = append(xssDetect, compile)
	}
	log.Println("xss攻击防护规则加载成功")
	// 添加xss攻击检测中间件
	r.Use(middleware.XSSDetectMiddleware(xssDetect))
	log.Println("xss攻击防护中间件加载成功")

	// 设置反向代理
	target, _ := url.Parse(conf.Server.TargetAddress)
	proxy := httputil.NewSingleHostReverseProxy(target)
	r.NoRoute(func(c *gin.Context) {
		proxy.ServeHTTP(c.Writer, c.Request)
	})
}

func RegisterUserHandler(r *gin.RouterGroup, dbs *Databases, conf *config.Config) {
	userService := service.NewAdminService(dbs.adminRepository)
	userController := controller.NewAdminController(userService, conf)
	authGroup := r.Group("/auth")
	authGroup.POST("/dologin", userController.DoLogin)
	authGroup.GET("/getCaptcha", userController.GetCaptcha)
}

func RegisterIPHandler(r *gin.RouterGroup, dbs *Databases, conf *config.Config) {
	ipService := service.NewIPListService(dbs.ipRepository)
	ipController := controller.NewIPController(ipService)
	ipGroup := r.Group("/ip")
	ipGroup.Use(middleware.WafAPIAuthMiddleware(conf.Secret.JwtSecretKey, dbs.adminRepository))
	ipGroup.POST("", ipController.CreateIP)
	ipGroup.GET("", ipController.FindPaginatedIP)
	ipGroup.PATCH(":id", ipController.UpdateIP)
	ipGroup.DELETE(":id", ipController.DeleteIP)
}

func RegisterRoutingHandler(r *gin.RouterGroup, dbs *Databases, conf *config.Config) {
	routingService := service.NewRoutingService(dbs.routingRepository)
	routingController := controller.NewRoutingController(routingService)
	routingGroup := r.Group("/routing")
	routingGroup.Use(middleware.WafAPIAuthMiddleware(conf.Secret.JwtSecretKey, dbs.adminRepository))
	routingGroup.POST("", routingController.CreateRouting)
	routingGroup.GET("", routingController.FindPaginatedRouting)
	routingGroup.PATCH(":id", routingController.UpdateRouting)
	routingGroup.DELETE(":id", routingController.DeleteRouting)
}

func RegisterLogHandler(r *gin.RouterGroup, dbs *Databases, conf *config.Config) {
	logService := service.NewLogService(dbs.logRepository)
	logController := controller.NewLogController(logService)
	logGroup := r.Group("/log")
	logGroup.Use(middleware.WafAPIAuthMiddleware(conf.Secret.JwtSecretKey, dbs.adminRepository))
	logGroup.GET("", logController.FindLogs)
	logGroup.GET("/getBlockLog", logController.FindPaginatedBlockedLog)
}

func RegisterConfigHandler(r *gin.RouterGroup, dbs *Databases, conf *config.Config) {
	configController := controller.NewConfigController(conf)
	configGroup := r.Group("/config")
	configGroup.Use(middleware.WafAPIAuthMiddleware(conf.Secret.JwtSecretKey, dbs.adminRepository))
	configGroup.GET("", configController.GetConfig)
	configGroup.PUT("", configController.UpdateConfig)
}

func RegisterSqlInjectHandler(r *gin.RouterGroup, dbs *Databases, conf *config.Config) {
	sqlInjectService := service.NewSqlInjectService(dbs.sqlInjectRepository)
	sqlInjectController := controller.NewSqlInjectController(sqlInjectService)
	sqlInjectGroup := r.Group("/sqlInject")
	sqlInjectGroup.Use(middleware.WafAPIAuthMiddleware(conf.Secret.JwtSecretKey, dbs.adminRepository))
	sqlInjectGroup.GET("", sqlInjectController.FindPaginatedSqlInject)
	sqlInjectGroup.POST("", sqlInjectController.CreateRule)
	sqlInjectGroup.PATCH(":id", sqlInjectController.UpdateRule)
	sqlInjectGroup.DELETE(":id", sqlInjectController.DeleteRule)
}

func RegisterXssDetectHandler(r *gin.RouterGroup, dbs *Databases, conf *config.Config) {
	xssDetectService := service.NewXssDetectService(dbs.xssDetectRepository)
	xssDetectController := controller.NewXssDetectController(xssDetectService)
	xssDetectGroup := r.Group("/xssDetect")
	xssDetectGroup.Use(middleware.WafAPIAuthMiddleware(conf.Secret.JwtSecretKey, dbs.adminRepository))
	xssDetectGroup.GET("", xssDetectController.FindPaginatedSqlInject)
	xssDetectGroup.POST("", xssDetectController.CreateRule)
	xssDetectGroup.PATCH(":id", xssDetectController.UpdateRule)
	xssDetectGroup.DELETE(":id", xssDetectController.DeleteRule)
}
