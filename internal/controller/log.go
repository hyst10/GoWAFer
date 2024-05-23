package controller

import (
	"GoWAFer/internal/service"
	"GoWAFer/pkg/pagination"
	"GoWAFer/pkg/utils/api_helper"
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
// @Tags Log（日志模块）
// @Produce json
// @Param days query string false "查询范围天数"
// @Success 200 {object} api_helper.Response
// @Router /waf/api/v1/log [get]
func (c *LogController) FindLogs(g *gin.Context) {
	days := g.DefaultQuery("days", "0.5")
	dayFloat, err := strconv.ParseFloat(days, 64)
	if err != nil {
		api_helper.ClientErrorHandler(g, 40004)
		return
	}
	var hours int
	if dayFloat < 1 {
		hours = 2
	}
	items := c.logService.FindLogs(int(dayFloat), hours)
	g.JSON(http.StatusOK, api_helper.Response{Status: 0, Msg: "success", Data: items})
}

// FindPaginatedBlockedLog godoc
// @Summary 分页查询被拦截的流量日志
// @Description 分页查询被拦截的流量日志
// @Tags Log（日志模块）
// @Produce json
// @Param ip query string false "ip搜索"
// @Param block_by query string false "封禁原因搜索"
// @Param method query string false "HTTP方法搜索"
// @Param orderDir query string false "日期排序"
// @Param page query int false "页码"
// @Param perPage query int false "页面大小"
// @Success 200 {object} api_helper.Response
// @Router /waf/api/v1/log/getBlockLog [get]
func (c *LogController) FindPaginatedBlockedLog(g *gin.Context) {
	ip := g.Query("ip")
	blockBy := g.Query("block_by")
	method := g.Query("method")
	order := g.DefaultQuery("orderDir", "desc")

	query := map[string]interface{}{
		"ip":      ip,
		"BlockBy": blockBy,
		"Method":  method,
	}

	page := pagination.NewFromRequest(g)
	page = c.logService.FindPaginatedLogs(page, query, order)
	g.JSON(http.StatusOK, api_helper.Response{Status: 0, Msg: "success", Data: page})
}
