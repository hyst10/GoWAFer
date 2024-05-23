package api

import (
	"GoWAFer/constants"
	"GoWAFer/internal/controller"
	"GoWAFer/internal/middleware"
	"GoWAFer/internal/repository"
	"GoWAFer/internal/service"
	"GoWAFer/pkg/config"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
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
	RegisterIPHandlers(apiGroup, dbs, conf)
	RegisterRoutingManageHandlers(apiGroup, dbs, conf)
	RegisterLogHandlers(apiGroup, dbs, conf)
	RegisterConfigHandler(apiGroup, dbs, conf)
	RegisterSqlInjectHandler(apiGroup, dbs, conf)
	RegisterXssDetectHandler(apiGroup, dbs, conf)

	// 添加日志中间件
	r.Use(middleware.TrafficLogger(dbs.logRepository))
	fmt.Println("日志中间件加载成功")

	// 添加IP管理中间件
	r.Use(middleware.IPManager(dbs.ipManageRepository))
	fmt.Println("IP管理中间件加载成功")

	// 添加路由守卫中间件
	r.Use(middleware.PathManager(dbs.routingManageRepository))
	fmt.Println("路径管理中间件加载成功")

	// 添加限速器中间件
	r.Use(middleware.RateLimitMiddleware(conf, dbs.ipManageRepository))
	fmt.Println("CC攻击防护中间件加载成功")

	// 添加CSRFToken中间件
	//r.Use(middleware.CsrfTokenMiddleware())
	//fmt.Println("CSRFToken中间件加载成功")

	// 加载sql注入防护规则
	sqlInjectRules, _ := dbs.sqlInjectRepository.GetAll()
	if len(sqlInjectRules) == 0 {
		// 导入规则
		ctx := context.Background()
		for rule := range constants.SqlInjectRules {
			rdb.SAdd(ctx, constants.SqlInjectKey, rule.String())
		}
	} else {
		// 使用redis中的规则
		constants.SqlInjectRules = make(map[*regexp.Regexp]bool)
		ctx := context.Background()
		for _, rule := range sqlInjectRules {
			compile := regexp.MustCompile(rule.Rule)
			constants.SqlInjectRules[compile] = true
			rdb.SAdd(ctx, constants.SqlInjectKey, rule)
		}
	}
	fmt.Println("sql注入防护规则加载成功")

	// 加载xss攻击防护规则
	xssDetectRules, _ := dbs.xssDetectRepository.GetAll()
	if len(xssDetectRules) == 0 {
		// 导入规则
		ctx := context.Background()
		for rule := range constants.XssDetectRules {
			rdb.SAdd(ctx, constants.XssDetectKey, rule.String())
		}
	} else {
		// 使用redis中的规则
		constants.XssDetectRules = make(map[*regexp.Regexp]bool)
		ctx := context.Background()
		for _, rule := range xssDetectRules {
			compile := regexp.MustCompile(rule.Rule)
			constants.XssDetectRules[compile] = true
			rdb.SAdd(ctx, constants.XssDetectKey, rule)
		}
	}
	fmt.Println("xss攻击防护规则加载成功")

	// 添加安全检测中间件
	r.Use(middleware.SecureRequestMiddleware())
	fmt.Println("安全检测中间件加载成功")

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

// RegisterIPHandlers 注册IP管理方法
func RegisterIPHandlers(r *gin.RouterGroup, dbs *Databases, conf *config.Config) {
	ipManageService := service.NewIPManageService(dbs.ipManageRepository)
	ipManageController := controller.NewIPManageController(ipManageService)
	ipManageGroup := r.Group("/ip")
	ipManageGroup.Use(middleware.WafAPIAuthMiddleware(conf.Secret.JwtSecretKey, dbs.adminRepository))
	ipManageGroup.POST("", ipManageController.AddIP)
	ipManageGroup.GET("", ipManageController.FindPaginatedIP)
	ipManageGroup.DELETE("", ipManageController.DeleteIP)
}

// RegisterRoutingManageHandlers 注册路由管理方法
func RegisterRoutingManageHandlers(r *gin.RouterGroup, dbs *Databases, conf *config.Config) {
	routingManageService := service.NewRoutingManageService(dbs.routingManageRepository)
	routingManageController := controller.NewRoutingManageController(routingManageService)
	routingGroup := r.Group("/routing")
	routingGroup.Use(middleware.WafAPIAuthMiddleware(conf.Secret.JwtSecretKey, dbs.adminRepository))
	routingGroup.POST("", routingManageController.AddRouting)
	routingGroup.GET("", routingManageController.GetRouting)
	routingGroup.DELETE("", routingManageController.DeleteRouting)
}

// RegisterLogHandlers 注册日志路由
func RegisterLogHandlers(r *gin.RouterGroup, dbs *Databases, conf *config.Config) {
	logService := service.NewLogService(dbs.logRepository)
	logController := controller.NewLogController(logService)
	logGroup := r.Group("/log")
	logGroup.Use(middleware.WafAPIAuthMiddleware(conf.Secret.JwtSecretKey, dbs.adminRepository))
	logGroup.GET("", logController.FindLogs)
	logGroup.GET("/getBlockLog", logController.FindPaginatedBlockedLog)
}

// RegisterConfigHandler 注册配置方法
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
