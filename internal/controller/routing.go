package controller

import (
	"GoWAFer/internal/model"
	"GoWAFer/internal/service"
	"GoWAFer/pkg/pagination"
	"GoWAFer/pkg/utils/api_handler"
	"github.com/gin-gonic/gin"
	"net/http"
	"regexp"
)

type RoutingController struct {
	routingService *service.RoutingService
}

func NewRoutingController(routingService *service.RoutingService) *RoutingController {
	return &RoutingController{routingService: routingService}
}

// CreateRouting godoc
// @Summary 新增路由
// @Description 新增路由
// @Tags Routing
// @Accept json
// @Product json
// @Param api_handler.CreateRoutingRequest body api_handler.CreateRoutingRequest true "请求体"
// @Success 200 {object} api_handler.Response
// @Router /waf/api/v1/routing [post]
func (c *RoutingController) CreateRouting(g *gin.Context) {
	var req api_handler.CreateRoutingRequest
	if err := g.ShouldBindJSON(&req); err != nil {
		api_handler.ClientErrorHandler(g, 40001)
		return
	}

	// 检查路由格式是否正确
	pattern := `^/[A-Za-z0-9/_]*$`
	matched, _ := regexp.MatchString(pattern, req.Route)
	if !matched {
		api_handler.ClientErrorHandler(g, 40006)
		return
	}

	newRouting := model.Routing{Type: req.Type, Route: req.Route, Method: req.Method}

	err := c.routingService.CreateRouting(&newRouting)
	if err != nil {
		api_handler.InternalErrorHandler(g, err)
		return
	}

	g.JSON(http.StatusOK, api_handler.Response{Status: 0, Message: "操作成功！"})
}

// FindPaginatedRouting godoc
// @Summary 分页查询路由
// @Description 分页查询路由
// @Tags Routing
// @Produce json
// @Param keywords query string false "关键词"
// @Param type query string true "路由类型"
// @Param page query int false "页码"
// @Param perPage query int false "页面大小"
// @Success 200 {object} api_handler.Response
// @Router /waf/api/v1/routing [get]
func (c *RoutingController) FindPaginatedRouting(g *gin.Context) {
	// 通过请求实例化分页结构体
	page := pagination.NewFromRequest(g)
	keyword := g.Query("keywords")
	routerType := g.Query("type")
	page = c.routingService.FindPaginatedRouters(page, routerType, keyword)
	g.JSON(http.StatusOK, api_handler.Response{Status: 0, Message: "success", Data: page})
}

// UpdateRouting godoc
// @Summary 编辑路由
// @Description 编辑路由
// @Tags Routing
// @Accept json
// @Produce json
// @Param id path int true "主键ID"
// @Param api_handler.UpdateRoutingRequest body api_handler.UpdateRoutingRequest true "请求体"
// @Success 200 {object} api_handler.Response
// @Router /waf/api/v1/routing/{id} [patch]
func (c *RoutingController) UpdateRouting(g *gin.Context) {
	var req api_handler.UpdateRoutingRequest
	if err := g.ShouldBindJSON(&req); err != nil {
		api_handler.ClientErrorHandler(g, 40001)
		return
	}

	// 检查路由格式是否正确
	pattern := `^/[A-Za-z0-9/_]*$`
	matched, _ := regexp.MatchString(pattern, req.Route)
	if !matched {
		api_handler.ClientErrorHandler(g, 40006)
		return
	}

	id, err := api_handler.GetUintParamFromPath(g, "id")
	if err != nil {
		api_handler.ClientErrorHandler(g, 40004)
		return
	}
	current, err := c.routingService.FindRoutingByID(id)
	if err != nil {
		api_handler.ClientErrorHandler(g, 40005)
		return
	}

	current.Route = req.Route
	current.Method = req.Method

	err = c.routingService.UpdateRouting(current)
	if err != nil {
		api_handler.InternalErrorHandler(g, err)
		return
	}

	g.JSON(http.StatusOK, api_handler.Response{Status: 0, Message: "操作成功！"})
}

// DeleteRouting godoc
// @Summary 删除路由
// @Description 删除路由
// @Tags Routing
// @Produce json
// @Param id path int true "主键ID"
// @Success 200 {object} api_handler.Response
// @Router /waf/api/v1/routing/{id} [delete]
func (c *RoutingController) DeleteRouting(g *gin.Context) {
	id, err := api_handler.GetUintParamFromPath(g, "id")
	if err != nil {
		api_handler.ClientErrorHandler(g, 40004)
		return
	}
	current, err := c.routingService.FindRoutingByID(id)
	if err != nil {
		api_handler.ClientErrorHandler(g, 40005)
		return
	}
	err = c.routingService.DeleteRouting(current)
	if err != nil {
		api_handler.InternalErrorHandler(g, err)
		return
	}
	g.JSON(http.StatusOK, api_handler.Response{Status: 0, Message: "操作成功！"})
}
