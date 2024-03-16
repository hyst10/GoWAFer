package controller

import (
	"GoWAFer/internal/service"
	"GoWAFer/pkg/pagination"
	"GoWAFer/pkg/utils/api_handler"
	"github.com/gin-gonic/gin"
	"net/http"
)

type BlockLogController struct {
	blockLogService *service.BlockLogService
}

func NewBlockLogController(blockLogService *service.BlockLogService) *BlockLogController {
	return &BlockLogController{blockLogService: blockLogService}
}

// FindPaginatedBlockLog godoc
// @Summary 分页查询拦截日志
// @Description 分页查询拦截日志
// @Tags BlockLog
// @Produce json
// @Param keywords query string false "查询IP"
// @Param page query int false "页码"
// @Param perPage query int false "页面大小"
// @Success 200 {object} api_handler.Response
// @Router /waf/api/v1/blockLog [get]
func (c *BlockLogController) FindPaginatedBlockLog(g *gin.Context) {
	// 通过请求实例化分页结构体
	page := pagination.NewFromRequest(g)
	keyword := g.Query("keywords")
	page = c.blockLogService.FindPaginatedLogs(page, keyword)
	g.JSON(http.StatusOK, api_handler.Response{Status: 0, Message: "success", Data: page})
}
