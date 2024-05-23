package controller

import (
	"GoWAFer/constants"
	"GoWAFer/internal/service"
	"GoWAFer/internal/types"
	"GoWAFer/pkg/utils/api_helper"
	"github.com/gin-gonic/gin"
	"net/http"
	"regexp"
)

type SqlInjectController struct {
	sqlInjectService *service.SqlInjectService
}

func NewSqlInjectController(sqlInjectService *service.SqlInjectService) *SqlInjectController {
	return &SqlInjectController{sqlInjectService: sqlInjectService}
}

// CreateRule godoc
// @Summary 新增sql注入防护规则
// @Description 新增sql注入防护规则
// @Tags SqlInject
// @Accept json
// @Product json
// @Param types.AddSqlInjectRuleRequest body types.AddSqlInjectRuleRequest true "请求体"
// @Success 200 {object} api_helper.Response
// @Router /waf/api/v1/sqlInject [post]
func (c *SqlInjectController) CreateRule(g *gin.Context) {
	var req types.AddSqlInjectRuleRequest
	if err := g.ShouldBindJSON(&req); err != nil {
		g.Error(types.ErrInvalidBody)
		return
	}

	err := c.sqlInjectService.AddRule(req.Rule)
	if err != nil {
		g.Error(err)
		return
	}

	g.JSON(http.StatusOK, api_helper.Response{Status: 0, Msg: "操作成功！"})

	go func() {
		constants.SqlInjectRules[regexp.MustCompile(req.Rule)] = true
	}()
}

// GetAllSqlInjectRules godoc
// @Summary 查询全部sql注入防护规则
// @Description 查询全部sql注入防护规则
// @Tags SqlInject
// @Produce json
// @Success 200 {object} api_helper.Response
// @Router /waf/api/v1/sqlInject [get]
func (c *SqlInjectController) GetAllSqlInjectRules(g *gin.Context) {
	items, count := c.sqlInjectService.GetAllRules()
	data := map[string]interface{}{
		"item":  items,
		"count": count,
	}
	g.JSON(http.StatusOK, api_helper.Response{Status: 0, Msg: "success", Data: data})
}

// DeleteRule godoc
// @Summary 删除sql注入防护规则
// @Description 删除sql注入防护规则
// @Tags SqlInject
// @Accept json
// @Product json
// @Param types.DeleteSqlInjectRuleRequest body types.DeleteSqlInjectRuleRequest true "请求体"
// @Success 200 {object} api_helper.Response
// @Router /waf/api/v1/sqlInject [delete]
func (c *SqlInjectController) DeleteRule(g *gin.Context) {
	var req types.DeleteSqlInjectRuleRequest
	if err := g.ShouldBindJSON(&req); err != nil {
		g.Error(types.ErrInvalidBody)
		return
	}

	err := c.sqlInjectService.DeleteRule(req.Rule)
	if err != nil {
		g.Error(err)
		return
	}

	g.JSON(http.StatusOK, api_helper.Response{Status: 0, Msg: "操作成功！"})

	go func() {
		delete(constants.SqlInjectRules, regexp.MustCompile(req.Rule))
	}()
}
