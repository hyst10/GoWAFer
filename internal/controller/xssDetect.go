package controller

import (
	"GoWAFer/internal/service"
	"GoWAFer/internal/types"
	"GoWAFer/pkg/utils/api_helper"
	"github.com/gin-gonic/gin"
	"net/http"
	"regexp"
)

type XssDetectController struct {
	xssDetectService *service.XssDetectService
}

func NewXssDetectController(xssDetectService *service.XssDetectService) *XssDetectController {
	return &XssDetectController{xssDetectService: xssDetectService}
}

// CreateRule godoc
// @Summary 新增xss攻击防护规则
// @Description 新增sql注入防护规则
// @Tags XssDetect
// @Accept json
// @Product json
// @Param types.XssDetectRule body types.XssDetectRule true "请求体"
// @Success 200 {object} api_helper.Response
// @Router /waf/api/v1/xssDetect [post]
func (c *XssDetectController) CreateRule(g *gin.Context) {
	var req types.XssDetectRule
	if err := g.ShouldBindJSON(&req); err != nil {
		g.Error(types.ErrInvalidBody)
		return
	}

	err := c.xssDetectService.AddRule(req.Rule)
	if err != nil {
		g.Error(err)
		return
	}

	g.JSON(http.StatusOK, api_helper.Response{Status: 0, Msg: "操作成功！"})

	go func() {
		types.XssDetectRules[regexp.MustCompile(req.Rule)] = true
	}()
}

// GetAllRules godoc
// @Summary 查询全部xss攻击防护规则
// @Description 查询全部xss攻击防护规则
// @Tags XssDetect
// @Produce json
// @Success 200 {object} api_helper.Response
// @Router /waf/api/v1/xssDetect [get]
func (c *XssDetectController) GetAllRules(g *gin.Context) {
	items, count := c.xssDetectService.GetAllRules()
	data := map[string]interface{}{
		"item":  items,
		"count": count,
	}
	g.JSON(http.StatusOK, api_helper.Response{Status: 0, Msg: "success", Data: data})
}

// DeleteRule godoc
// @Summary 删除xss攻击防护规则
// @Description 删除xss攻击防护规则
// @Tags XssDetect
// @Accept json
// @Product json
// @Param types.DeleteXssDetectRuleRequest body types.DeleteXssDetectRuleRequest true "请求体"
// @Success 200 {object} api_helper.Response
// @Router /waf/api/v1/xssDetect [delete]
func (c *XssDetectController) DeleteRule(g *gin.Context) {
	var req types.DeleteXssDetectRuleRequest
	if err := g.ShouldBindJSON(&req); err != nil {
		g.Error(types.ErrInvalidBody)
		return
	}

	err := c.xssDetectService.DeleteRule(req.Rule)
	if err != nil {
		g.Error(err)
		return
	}

	g.JSON(http.StatusOK, api_helper.Response{Status: 0, Msg: "操作成功！"})

	go func() {
		delete(types.XssDetectRules, regexp.MustCompile(req.Rule))
	}()
}
