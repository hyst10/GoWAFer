package web

import (
	"GoWAFer/pkg/config"
	"github.com/gin-gonic/gin"
	"net/http"
)

func RegisterWebHandler(r *gin.Engine, conf *config.Config) {
	wafGroup := r.Group("/waf")
	wafGroup.GET("/login", func(c *gin.Context) {
		c.HTML(http.StatusOK, "login.html", nil)
	})
}
