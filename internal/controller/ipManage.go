package controller

import (
	"GoWAFer/internal/service"
	"GoWAFer/internal/types"
	"GoWAFer/pkg/pagination"
	"GoWAFer/pkg/utils"
	"GoWAFer/pkg/utils/api_helper"
	"github.com/gin-gonic/gin"
	"net/http"
)

type IPManageController struct {
	ipManageService *service.IPManageService
}

func NewIPManageController(s *service.IPManageService) *IPManageController {
	return &IPManageController{ipManageService: s}
}

// AddIP godoc
// @Summary 新增IP管理记录
// @Description 新增IP记录
// @Tags IP
// @Accept json
// @Product json
// @Param types.AddIPRequest body types.AddIPRequest true "请求体"
// @Success 200 {object} api_helper.Response
// @Router /waf/api/v1/ip [post]
func (c *IPManageController) AddIP(g *gin.Context) {
	var req types.AddIPRequest
	if err := g.ShouldBindJSON(&req); err != nil {
		g.Error(types.ErrInvalidBody)
		return
	}

	// 检查IP格式是否正确
	if !utils.ValidateIP(req.IP) {
		g.Error(types.ErrIPNotMatch)
		return
	}

	// 新增IP记录
	err := c.ipManageService.AddIP(req.IP, req.Expiration, req.IsBlack)
	if err != nil {
		g.Error(err)
		return
	}

	g.JSON(http.StatusOK, api_helper.Response{Status: 0, Msg: "操作成功！"})
}

// FindPaginatedIP godoc
// @Summary 分页查询IP管理记录
// @Description 分页查询IP管理记录
// @Tags IP
// @Produce json
// @Param query query string false "查询关键字"
// @Param isBlack query string false "是否为黑名单"
// @Param page query int false "页码"
// @Param perPage query int false "页面大小"
// @Success 200 {object} api_helper.Response
// @Router /waf/api/v1/ip [get]
func (c *IPManageController) FindPaginatedIP(g *gin.Context) {
	// 通过请求实例化分页结构体
	page := pagination.NewFromRequest(g)
	query := g.Query("query")
	isBlack := g.DefaultQuery("isBlack", "false")
	var isBlackBool bool
	if isBlack == "true" {
		isBlackBool = true
	} else {
		isBlackBool = false
	}
	page = c.ipManageService.GetIPWithPagination(page, isBlackBool, query)
	g.JSON(http.StatusOK, api_helper.Response{Status: 0, Msg: "success", Data: page})
}

// DeleteIP godoc
// @Summary 删除IP记录
// @Description 删除IP记录
// @Tags IP
// @Accept json
// @Product json
// @Param types.DeleteIPRequest body types.DeleteIPRequest true "请求体"
// @Success 200 {object} api_helper.Response
// @Router /waf/api/v1/ip [delete]
func (c *IPManageController) DeleteIP(g *gin.Context) {
	var req types.DeleteIPRequest
	if err := g.ShouldBindJSON(&req); err != nil {
		g.Error(types.ErrInvalidBody)
		return
	}

	err := c.ipManageService.DeleteIP(req.IP, req.IsBlack)
	if err != nil {
		g.Error(err)
		return
	}

	g.JSON(http.StatusOK, api_helper.Response{Status: 0, Msg: "操作成功！"})
}
