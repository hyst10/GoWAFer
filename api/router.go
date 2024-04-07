package api

import (
	"GoWAFer/internal/controller"
	"GoWAFer/internal/middleware"
	"GoWAFer/internal/repository"
	"GoWAFer/internal/service"
	"GoWAFer/internal/types"
	"GoWAFer/pkg/config"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
	"log"
	"net/http/httputil"
	"net/url"
	"regexp"
)

type Databases struct {
	adminRepository         *repository.AdminRepository
	logRepository           *repository.LogRepository
	ipManageRepository      *repository.IPManageRepository
	routingManageRepository *repository.RoutingManageRepository
	sqlInjectRepository     *repository.SqlInjectRepository
	xssDetectRepository     *repository.XssDetectRepository
}

func NewDatabases(db *gorm.DB, rdb *redis.Client) *Databases {
	return &Databases{
		adminRepository:         repository.NewAdminRepository(db),
		logRepository:           repository.NewLogRepository(db),
		ipManageRepository:      repository.NewIPManageRepository(rdb),
		routingManageRepository: repository.NewRoutingManageRepository(rdb),
		sqlInjectRepository:     repository.NewSqlInjectRepository(rdb),
		xssDetectRepository:     repository.NewXssDetectRepository(rdb),
	}
}

func RegisterAllHandlers(r *gin.Engine, db *gorm.DB, rdb *redis.Client, conf *config.Config) {
	dbs := NewDatabases(db, rdb)

	apiGroup := r.Group("/waf/api/v1")
	apiGroup.Use(middleware.ErrorHandlingMiddleware()) // 错误处理中间件

	RegisterUserHandler(apiGroup, dbs, conf)
	RegisterIPHandler(apiGroup, dbs, conf)
	RegisterRoutingManageHandler(apiGroup, dbs, conf)
	RegisterLogHandler(apiGroup, dbs, conf)
	RegisterConfigHandler(apiGroup, dbs, conf)
	RegisterSqlInjectHandler(apiGroup, dbs, conf)
	RegisterXssDetectHandler(apiGroup, dbs, conf)

	// 添加日志中间件
	r.Use(middleware.TrafficLogger(dbs.logRepository))
	log.Println("日志中间件加载成功")
	// 添加IP管理中间件
	r.Use(middleware.IPManager(dbs.ipManageRepository))
	log.Println("IP管理中间件加载成功")
	// 添加路由守卫中间件
	r.Use(middleware.RouteGuardMiddleware(dbs.routingManageRepository))
	log.Println("路由守卫中间件加载成功")
	// 添加限速器中间件
	r.Use(middleware.RateLimitMiddleware(conf, dbs.ipManageRepository))
	log.Println("CC攻击防护中间件加载成功")
	// 添加CSRFToken中间件
	r.Use(middleware.CsrfTokenMiddleware())
	log.Println("CSRFToken中间件加载成功")
	// 加载sql注入防护规则
	sqlInjectRules, _ := dbs.sqlInjectRepository.GetAll()
	if len(sqlInjectRules) == 0 {
		// 导入规则
		ctx := context.Background()
		for rule := range types.SqlInjectRules {
			rdb.SAdd(ctx, repository.SqlInjectRules, rule.String())
		}
	} else {
		// 使用redis中的规则
		types.SqlInjectRules = make(map[*regexp.Regexp]bool)
		ctx := context.Background()
		for _, rule := range sqlInjectRules {
			compile := regexp.MustCompile(rule.Rule)
			types.SqlInjectRules[compile] = true
			rdb.SAdd(ctx, repository.SqlInjectRules, rule)
		}
	}
	log.Println("sql注入防护规则加载成功")
	// 添加sql注入检测中间件
	r.Use(middleware.SqlInjectMiddleware())
	log.Println("sql注入防护中间件加载成功")
	// 加载xss攻击防护规则
	xssDetectRules, _ := dbs.xssDetectRepository.GetAll()
	if len(xssDetectRules) == 0 {
		// 导入规则
		ctx := context.Background()
		for rule := range types.XssDetectRules {
			rdb.SAdd(ctx, repository.XssDetectRules, rule.String())
		}
	} else {
		// 使用redis中的规则
		types.XssDetectRules = make(map[*regexp.Regexp]bool)
		ctx := context.Background()
		for _, rule := range xssDetectRules {
			compile := regexp.MustCompile(rule.Rule)
			types.XssDetectRules[compile] = true
			rdb.SAdd(ctx, repository.XssDetectRules, rule)
		}
	}
	log.Println("xss攻击防护规则加载成功")
	// 添加xss攻击检测中间件
	r.Use(middleware.XSSDetectMiddleware())
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

// RegisterIPHandler 注册IP管理方法
func RegisterIPHandler(r *gin.RouterGroup, dbs *Databases, conf *config.Config) {
	ipManageService := service.NewIPManageService(dbs.ipManageRepository)
	ipManageController := controller.NewIPManageController(ipManageService)
	ipManageGroup := r.Group("/ip")
	ipManageGroup.Use(middleware.WafAPIAuthMiddleware(conf.Secret.JwtSecretKey, dbs.adminRepository))
	ipManageGroup.POST("", ipManageController.AddIP)
	ipManageGroup.GET("", ipManageController.FindPaginatedIP)
	ipManageGroup.DELETE("", ipManageController.DeleteIP)
}

// RegisterRoutingManageHandler 注册路由管理方法
func RegisterRoutingManageHandler(r *gin.RouterGroup, dbs *Databases, conf *config.Config) {
	routingManageService := service.NewRoutingManageService(dbs.routingManageRepository)
	routingManageController := controller.NewRoutingManageController(routingManageService)
	routingGroup := r.Group("/routing")
	routingGroup.Use(middleware.WafAPIAuthMiddleware(conf.Secret.JwtSecretKey, dbs.adminRepository))
	routingGroup.POST("", routingManageController.AddRouting)
	routingGroup.GET("", routingManageController.GetRouting)
	routingGroup.DELETE("", routingManageController.DeleteRouting)
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

// RegisterSqlInjectHandler 注册sql注入防护方法
func RegisterSqlInjectHandler(r *gin.RouterGroup, dbs *Databases, conf *config.Config) {
	sqlInjectService := service.NewSqlInjectService(dbs.sqlInjectRepository)
	sqlInjectController := controller.NewSqlInjectController(sqlInjectService)
	sqlInjectGroup := r.Group("/sqlInject")
	sqlInjectGroup.Use(middleware.WafAPIAuthMiddleware(conf.Secret.JwtSecretKey, dbs.adminRepository))
	sqlInjectGroup.GET("", sqlInjectController.GetAllSqlInjectRules)
	sqlInjectGroup.POST("", sqlInjectController.CreateRule)
	sqlInjectGroup.DELETE("", sqlInjectController.DeleteRule)
}

// RegisterXssDetectHandler 注册xss攻击防护方法
func RegisterXssDetectHandler(r *gin.RouterGroup, dbs *Databases, conf *config.Config) {
	xssDetectService := service.NewXssDetectService(dbs.xssDetectRepository)
	xssDetectController := controller.NewXssDetectController(xssDetectService)
	xssDetectGroup := r.Group("/xssDetect")
	xssDetectGroup.Use(middleware.WafAPIAuthMiddleware(conf.Secret.JwtSecretKey, dbs.adminRepository))
	xssDetectGroup.GET("", xssDetectController.GetAllRules)
	xssDetectGroup.POST("", xssDetectController.CreateRule)
	xssDetectGroup.DELETE("", xssDetectController.DeleteRule)
}
