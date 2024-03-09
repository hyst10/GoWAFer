package controller

import (
	"GoWAFer/internal/service"
	"GoWAFer/pkg/utils/api_handler"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type LogController struct {
	logService *service.LogService
}

func NewLogController(logService *service.LogService) *LogController {
	return &LogController{
		logService: logService,
	}
}

// FindLogs godoc
// @Summary 查询指定天数和小时数的日志记录
// @Description 查询指定天数和小时数的日志记录
// @Tags Log
// @Produce json
// @Param days query string false "查询范围天数"
// @Success 200 {object} api_handler.Response
// @Router /waf/api/v1/log [get]
func (c *LogController) FindLogs(g *gin.Context) {
	days := g.DefaultQuery("days", "0.5")
	dayFloat, err := strconv.ParseFloat(days, 64)
	if err != nil {
		api_handler.ClientErrorHandler(g, 40004)
		return
	}
	var hours int
	if dayFloat < 1 {
		hours = 2
	}
	items := c.logService.FindLogs(int(dayFloat), hours)
	g.JSON(http.StatusOK, api_handler.Response{Status: 0, Message: "success", Data: items})
}
