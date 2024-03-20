package controller

import (
	"GoWAFer/pkg/config"
	"GoWAFer/pkg/utils/api_handler"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"net/http"
	"os"
	"time"
)

type ConfigController struct {
	conf *config.Config
}

func NewConfigController(conf *config.Config) *ConfigController {
	return &ConfigController{
		conf: conf,
	}
}

// GetConfig godoc
// @Summary 获取配置信息
// @Description 获取配置信息
// @Tags Config
// @Produce json
// @Success 200 {object} api_handler.Response
// @Router /waf/api/v1/config [get]
func (c *ConfigController) GetConfig(g *gin.Context) {
	g.JSON(http.StatusOK, api_handler.Response{Status: 0, Message: "success", Data: c.conf})
}

// UpdateConfig godoc
// @Summary 修改配置文件
// @Description 修改配置文件并重启服务
// @Tags Config
// @Produce json
// @Param config.Config body config.Config true "请求体"
// @Success 200 {object} api_handler.Response
// @Router /waf/api/v1/config [put]
func (c *ConfigController) UpdateConfig(g *gin.Context) {
	var req config.Config
	if err := g.ShouldBindJSON(&req); err != nil {
		api_handler.ClientErrorHandler(g, 40001)
		return
	}

	// 构造全量更新的配置map
	updateMap := map[string]interface{}{
		"server": map[string]interface{}{
			"targetAddress": req.Server.TargetAddress,
		},
		"secret": map[string]interface{}{
			"jwtSecretKey":     req.Secret.JwtSecretKey,
			"sessionSecretKey": req.Secret.SessionSecretKey,
		},
		"rateLimiter": map[string]interface{}{
			"maxCounter":  req.RateLimiter.MaxCounter,
			"banCounter":  req.RateLimiter.BanCounter,
			"banDuration": req.RateLimiter.BanDuration,
			"mode":        req.RateLimiter.Mode,
			"tokenBucket": map[string]interface{}{
				"maxToken":       req.RateLimiter.TokenBucket.MaxToken,
				"tokenPerSecond": req.RateLimiter.TokenBucket.TokenPerSecond,
			},
			"leakyBucket": map[string]interface{}{
				"capacity":       req.RateLimiter.LeakyBucket.Capacity,
				"leakyPerSecond": req.RateLimiter.LeakyBucket.LeakyPerSecond,
			},
			"fixedWindow": map[string]interface{}{
				"windowSize": req.RateLimiter.FixedWindow.WindowSize,
				"maxRequest": req.RateLimiter.FixedWindow.MaxRequest,
			},
			"slideWindow": map[string]interface{}{
				"windowSize": req.RateLimiter.SlideWindow.WindowSize,
				"maxRequest": req.RateLimiter.SlideWindow.MaxRequest,
			},
		},
	}

	// 使用Viper进行全量更新
	if err := viper.MergeConfigMap(updateMap); err != nil {
		api_handler.InternalErrorHandler(g, err)
		return
	}

	// 将更新后的配置保存到文件
	if err := viper.WriteConfig(); err != nil {
		api_handler.InternalErrorHandler(g, err)
		return
	}

	g.JSON(http.StatusOK, api_handler.Response{Status: 0, Message: "操作成功！"})

	// 启动一个新的goroutine来处理延迟重启逻辑
	go func() {
		// 等待足够的时间以确保响应已发送
		time.Sleep(3 * time.Second)
		// 退出程序
		os.Exit(1)
	}()
}
