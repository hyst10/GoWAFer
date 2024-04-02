package controller

import (
	"GoWAFer/internal/model"
	"GoWAFer/internal/service"
	"GoWAFer/pkg/pagination"
	"GoWAFer/pkg/utils/api_helper"
	"github.com/gin-gonic/gin"
	"net/http"
)

type SqlInjectController struct {
	sqlInjectService *service.SqlInjectService
}

func NewSqlInjectController(sqlInjectService *service.SqlInjectService) *SqlInjectController {
	return &SqlInjectController{sqlInjectService: sqlInjectService}
}

// CreateRule godoc
// @Summary 新增sql注入规则
// @Description 新增sql注入规则
// @Tags SQlInject
// @Accept json
// @Product json
// @Param api_helper.CreateSqlInjectRequest body api_helper.CreateSqlInjectRequest true "请求体"
// @Success 200 {object} api_helper.Response
// @Router /waf/api/v1/sqlInject [post]
func (c *SqlInjectController) CreateRule(g *gin.Context) {
	var req api_helper.CreateSqlInjectRequest
	if err := g.ShouldBindJSON(&req); err != nil {
		api_helper.ClientErrorHandler(g, 40001)
		return
	}

	newRule := model.SqlInjectionRules{Rule: req.Rule}

	err := c.sqlInjectService.CreateRule(&newRule)
	if err != nil {
		api_helper.InternalErrorHandler(g, err)
		return
	}

	g.JSON(http.StatusOK, api_helper.Response{Status: 0, Msg: "操作成功！"})
}

// FindPaginatedSqlInject godoc
// @Summary 分页查询SQL注入规则
// @Description 分页查询SQL注入规则
// @Tags SQlInject
// @Produce json
// @Param keywords query string false "关键词"
// @Param page query int false "页码"
// @Param perPage query int false "页面大小"
// @Success 200 {object} api_helper.Response
// @Router /waf/api/v1/sqlInject [get]
func (c *SqlInjectController) FindPaginatedSqlInject(g *gin.Context) {
	// 通过请求实例化分页结构体
	page := pagination.NewFromRequest(g)
	keyword := g.Query("keywords")
	page = c.sqlInjectService.FindPaginatedRules(page, keyword)
	g.JSON(http.StatusOK, api_helper.Response{Status: 0, Msg: "success", Data: page})
}

// UpdateRule godoc
// @Summary 编辑sql注入规则
// @Description 编辑sql注入规则
// @Tags SQlInject
// @Produce json
// @Param id path int true "主键ID"
// @Param api_helper.CreateSqlInjectRequest body api_helper.CreateSqlInjectRequest true "请求体"
// @Success 200 {object} api_helper.Response
// @Router /waf/api/v1/sqlInject/{id} [patch]
func (c *SqlInjectController) UpdateRule(g *gin.Context) {
	var req api_helper.CreateSqlInjectRequest
	if err := g.ShouldBindJSON(&req); err != nil {
		api_helper.ClientErrorHandler(g, 40001)
		return
	}

	id, err := api_helper.GetUintParamFromPath(g, "id")
	if err != nil {
		api_helper.ClientErrorHandler(g, 40004)
		return
	}
	current, err := c.sqlInjectService.FindRuleByID(id)
	if err != nil {
		api_helper.ClientErrorHandler(g, 40005)
		return
	}
	current.Rule = req.Rule
	err = c.sqlInjectService.UpdateRule(current)
	if err != nil {
		api_helper.InternalErrorHandler(g, err)
		return
	}

	g.JSON(http.StatusOK, api_helper.Response{Status: 0, Msg: "操作成功！"})
}

// DeleteRule godoc
// @Summary 删除sql注入规则
// @Description 删除sql注入规则
// @Tags SQlInject
// @Produce json
// @Param id path int true "主键ID"
// @Success 200 {object} api_helper.Response
// @Router /waf/api/v1/sqlInject/{id} [delete]
func (c *SqlInjectController) DeleteRule(g *gin.Context) {
	id, err := api_helper.GetUintParamFromPath(g, "id")
	if err != nil {
		api_helper.ClientErrorHandler(g, 40004)
		return
	}
	current, err := c.sqlInjectService.FindRuleByID(id)
	if err != nil {
		api_helper.ClientErrorHandler(g, 40005)
		return
	}
	err = c.sqlInjectService.DeleteRule(current)
	if err != nil {
		api_helper.InternalErrorHandler(g, err)
		return
	}

	g.JSON(http.StatusOK, api_helper.Response{Status: 0, Msg: "操作成功！"})
}
