package api

import (
	"GoWAFer/internal/controller"
	"GoWAFer/internal/middleware"
	"GoWAFer/internal/repository"
	"GoWAFer/internal/service"
	"GoWAFer/pkg/config"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
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
	apiGroup.Use(middleware.TrafficLogger(dbs.logRepository))
	RegisterUserHandler(apiGroup, dbs, conf)
	RegisterIPHandler(apiGroup, dbs, conf)
	RegisterLogHandler(apiGroup, dbs)
}

func RegisterUserHandler(r *gin.RouterGroup, dbs *Databases, conf *config.Config) {
	userService := service.NewUserService(dbs.userRepository)
	userController := controller.NewUserController(userService, conf)
	authGroup := r.Group("/auth")
	authGroup.POST("/dologin", userController.DoLogin)
	authGroup.GET("/getCaptcha", userController.GetCaptcha)
}

func RegisterIPHandler(r *gin.RouterGroup, dbs *Databases, conf *config.Config) {
	ipService := service.NewIPListService(dbs.ipRepository)
	ipController := controller.NewIPController(ipService)
	ipGroup := r.Group("/ip")
	ipGroup.POST("", ipController.CreateIP)
	ipGroup.GET("", ipController.FindPaginatedIP)
	ipGroup.PATCH(":id", ipController.UpdateIP)
	ipGroup.DELETE(":id", ipController.DeleteIP)
}

func RegisterLogHandler(r *gin.RouterGroup, dbs *Databases) {
	logService := service.NewLogService(dbs.logRepository)
	logController := controller.NewLogController(logService)
	logGroup := r.Group("/log")
	logGroup.GET("", logController.FindLogs)
}
