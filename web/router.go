package web

import (
	"GoWAFer/pkg/config"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func RegisterWebHandler(r *gin.Engine, conf *config.Config) {
	wafGroup := r.Group("/waf")
	// 登录页
	wafGroup.GET("/login", func(c *gin.Context) {
		c.HTML(http.StatusOK, "login.html", nil)
	})
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
