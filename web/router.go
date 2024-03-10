package web

import (
	"GoWAFer/internal/middleware"
	"GoWAFer/internal/repository"
	"GoWAFer/pkg/config"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
)

func RegisterWebHandler(r *gin.Engine, db *gorm.DB, conf *config.Config) {
	wafGroup := r.Group("/waf")

	// 登录页
	wafGroup.GET("/login", func(c *gin.Context) {
		c.HTML(http.StatusOK, "login.html", nil)
	})

	// 添加鉴权中间件
	userRepository := repository.NewUserRepository(db)
	wafGroup.Use(middleware.WafWebAuthMiddleware(conf.Secret.JwtSecretKey, userRepository))

	// 首页
	wafGroup.GET("", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})
	wafGroup.GET("/index", func(c *gin.Context) {
		c.Redirect(http.StatusFound, "/waf")
	})

	// json配置路由
	wafGroup.GET("/pages/:filename", func(c *gin.Context) {
		filename := fmt.Sprintf("%s.json", c.Param("filename"))
		jsonFilePath := fmt.Sprintf("./pages/%s", filename)
		// 设置响应头，指示浏览器不要缓存文件
		c.Header("Cache-Control", "no-cache, no-store, must-revalidate")
		c.File(jsonFilePath)
	})

	// APP多页应用
	wafGroup.GET("/app/*action", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

}
