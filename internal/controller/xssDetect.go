package controller

import (
	"GoWAFer/internal/model"
	"GoWAFer/internal/service"
	"GoWAFer/pkg/pagination"
	"GoWAFer/pkg/utils/api_handler"
	"github.com/gin-gonic/gin"
	"net/http"
)

type XssDetectController struct {
	xssDetectService *service.XssDetectService
}

func NewXssDetectController(xssDetectService *service.XssDetectService) *XssDetectController {
	return &XssDetectController{xssDetectService: xssDetectService}
}

// CreateRule godoc
// @Summary 新增xss攻击匹配规则
// @Description 新增xss攻击匹配规则
// @Tags XssDetect
// @Accept json
// @Product json
// @Param api_handler.CreateXssDetectRequest body api_handler.CreateXssDetectRequest true "请求体"
// @Success 200 {object} api_handler.Response
// @Router /waf/api/v1/xssDetect [post]
func (c *XssDetectController) CreateRule(g *gin.Context) {
	var req api_handler.CreateXssDetectRequest
	if err := g.ShouldBindJSON(&req); err != nil {
		api_handler.ClientErrorHandler(g, 40001)
		return
	}

	newRule := model.XssDetectRules{Rule: req.Rule}

	err := c.xssDetectService.CreateRule(&newRule)
	if err != nil {
		api_handler.InternalErrorHandler(g, err)
		return
	}

	g.JSON(http.StatusOK, api_handler.Response{Status: 0, Message: "操作成功！"})
}

// FindPaginatedSqlInject godoc
// @Summary 分页查询xss攻击匹配规则
// @Description 分页查询xss攻击匹配规则
// @Tags XssDetect
// @Produce json
// @Param keywords query string false "关键词"
// @Param page query int false "页码"
// @Param perPage query int false "页面大小"
// @Success 200 {object} api_handler.Response
// @Router /waf/api/v1/xssDetect [get]
func (c *XssDetectController) FindPaginatedSqlInject(g *gin.Context) {
	// 通过请求实例化分页结构体
	page := pagination.NewFromRequest(g)
	keyword := g.Query("keywords")
	page = c.xssDetectService.FindPaginatedRules(page, keyword)
	g.JSON(http.StatusOK, api_handler.Response{Status: 0, Message: "success", Data: page})
}

// UpdateRule godoc
// @Summary 编辑xss攻击匹配规则
// @Description 编辑xss攻击匹配规则
// @Tags XssDetect
// @Produce json
// @Param id path int true "主键ID"
// @Param api_handler.CreateXssDetectRequest body api_handler.CreateXssDetectRequest true "请求体"
// @Success 200 {object} api_handler.Response
// @Router /waf/api/v1/xssDetect/{id} [patch]
func (c *XssDetectController) UpdateRule(g *gin.Context) {
	var req api_handler.CreateXssDetectRequest
	if err := g.ShouldBindJSON(&req); err != nil {
		api_handler.ClientErrorHandler(g, 40001)
		return
	}

	id, err := api_handler.GetUintParamFromPath(g, "id")
	if err != nil {
		api_handler.ClientErrorHandler(g, 40004)
		return
	}
	current, err := c.xssDetectService.FindRuleByID(id)
	if err != nil {
		api_handler.ClientErrorHandler(g, 40005)
		return
	}
	current.Rule = req.Rule
	err = c.xssDetectService.UpdateRule(current)
	if err != nil {
		api_handler.InternalErrorHandler(g, err)
		return
	}

	g.JSON(http.StatusOK, api_handler.Response{Status: 0, Message: "操作成功！"})
}

// DeleteRule godoc
// @Summary 删除sql注入规则
// @Description 删除sql注入规则
// @Tags XssDetect
// @Produce json
// @Param id path int true "主键ID"
// @Success 200 {object} api_handler.Response
// @Router /waf/api/v1/xssDetect/{id} [delete]
func (c *XssDetectController) DeleteRule(g *gin.Context) {
	id, err := api_handler.GetUintParamFromPath(g, "id")
	if err != nil {
		api_handler.ClientErrorHandler(g, 40004)
		return
	}
	current, err := c.xssDetectService.FindRuleByID(id)
	if err != nil {
		api_handler.ClientErrorHandler(g, 40005)
		return
	}
	err = c.xssDetectService.DeleteRule(current)
	if err != nil {
		api_handler.InternalErrorHandler(g, err)
		return
	}

	g.JSON(http.StatusOK, api_handler.Response{Status: 0, Message: "操作成功！"})
}
