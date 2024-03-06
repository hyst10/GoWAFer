package api

import (
	"GoWAFer/internal/controller"
	"GoWAFer/internal/repository"
	"GoWAFer/internal/service"
	"GoWAFer/pkg/config"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Databases struct {
	userRepository *repository.UserRepository
}

func NewDatabases(db *gorm.DB) *Databases {
	return &Databases{
		userRepository: repository.NewUserRepository(db),
	}
}

func RegisterAllHandlers(r *gin.Engine, db *gorm.DB, conf *config.Config) {
	dbs := NewDatabases(db)
	apiGroup := r.Group("/waf/api/v1")
	RegisterUserHandler(apiGroup, dbs, conf)

}

func RegisterUserHandler(r *gin.RouterGroup, dbs *Databases, conf *config.Config) {
	userService := service.NewUserService(dbs.userRepository)
	userController := controller.NewUserController(userService, conf)
	authGroup := r.Group("/auth")
	authGroup.POST("/dologin", userController.DoLogin)
	authGroup.GET("/getCaptcha", userController.GetCaptcha)
}
