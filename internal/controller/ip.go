package controller

import (
	"GoWAFer/internal/model"
	"GoWAFer/internal/service"
	"GoWAFer/pkg/pagination"
	"GoWAFer/pkg/utils/api_helper"
	"github.com/gin-gonic/gin"
	"net"
	"net/http"
	"strconv"
	"time"
)

type IPController struct {
	ipService *service.IPService
}

func NewIPController(ipService *service.IPService) *IPController {
	return &IPController{ipService: ipService}
}

// CreateIP godoc
// @Summary 新增IP记录
// @Description 新增IP记录
// @Tags IP
// @Accept json
// @Product json
// @Param api_helper.CreateIPRequest body api_helper.CreateIPRequest true "请求体"
// @Success 200 {object} api_helper.Response
// @Router /waf/api/v1/ip [post]
func (c *IPController) CreateIP(g *gin.Context) {
	var req api_helper.CreateIPRequest
	if err := g.ShouldBindJSON(&req); err != nil {
		api_helper.ClientErrorHandler(g, 40001)
		return
	}

	// 检查IP格式是否正确
	if net.ParseIP(req.IPAddress) == nil {
		api_helper.ClientErrorHandler(g, 40007)
		return
	}

	expirationInt, err := strconv.Atoi(req.ExpirationAt)
	if err != nil {
		api_helper.ClientErrorHandler(g, 40001)
		return
	}

	expiration := time.Unix(int64(expirationInt), 0)
	newIP := model.IP{Type: req.Type, IPAddress: req.IPAddress, ExpirationAt: expiration}

	err = c.ipService.CreateIP(&newIP)
	if err != nil {
		api_helper.InternalErrorHandler(g, err)
		return
	}

	g.JSON(http.StatusOK, api_helper.Response{Status: 0, Msg: "操作成功！"})
}

// FindPaginatedIP godoc
// @Summary 分页查询IP
// @Description 分页查询IP
// @Tags IP
// @Produce json
// @Param keywords query string false "关键词"
// @Param type query string true "IP类型"
// @Param page query int false "页码"
// @Param perPage query int false "页面大小"
// @Success 200 {object} api_helper.Response
// @Router /waf/api/v1/ip [get]
func (c *IPController) FindPaginatedIP(g *gin.Context) {
	// 通过请求实例化分页结构体
	page := pagination.NewFromRequest(g)
	keyword := g.Query("keywords")
	ipType := g.Query("type")
	page = c.ipService.FindPaginatedIPs(page, ipType, keyword)
	g.JSON(http.StatusOK, api_helper.Response{Status: 0, Msg: "success", Data: page})
}

// UpdateIP godoc
// @Summary 编辑IP
// @Description 编辑IP
// @Tags IP
// @Produce json
// @Param id path int true "IP主键ID"
// @Param api_helper.UpdateIPRequest body api_helper.UpdateIPRequest true "请求体"
// @Success 200 {object} api_helper.Response
// @Router /waf/api/v1/ip/{id} [patch]
func (c *IPController) UpdateIP(g *gin.Context) {
	var req api_helper.UpdateIPRequest
	if err := g.ShouldBindJSON(&req); err != nil {
		api_helper.ClientErrorHandler(g, 40001)
		return
	}

	// 检查IP格式是否正确
	if net.ParseIP(req.IPAddress) == nil {
		api_helper.ClientErrorHandler(g, 40007)
		return
	}

	id, err := api_helper.GetUintParamFromPath(g, "id")
	if err != nil {
		api_helper.ClientErrorHandler(g, 40004)
		return
	}
	current, err := c.ipService.FindIPByID(id)
	if err != nil {
		api_helper.ClientErrorHandler(g, 40005)
		return
	}

	layout := "2006-01-02T15:04:05Z07:00"
	expiration, err := time.Parse(layout, req.ExpirationAt)
	if err != nil {
		// 尝试时间戳
		expirationInt, err := strconv.Atoi(req.ExpirationAt)
		if err != nil {
			api_helper.ClientErrorHandler(g, 40001)
			return
		}
		expiration = time.Unix(int64(expirationInt), 0)
	}

	current.IPAddress = req.IPAddress
	current.ExpirationAt = expiration

	err = c.ipService.UpdateIP(current)
	if err != nil {
		api_helper.InternalErrorHandler(g, err)
		return
	}

	g.JSON(http.StatusOK, api_helper.Response{Status: 0, Msg: "操作成功！"})
}

// DeleteIP godoc
// @Summary 删除IP
// @Description 删除IP
// @Tags IP
// @Produce json
// @Param id path int true "IP主键ID"
// @Success 200 {object} api_helper.Response
// @Router /waf/api/v1/ip/{id} [delete]
func (c *IPController) DeleteIP(g *gin.Context) {
	id, err := api_helper.GetUintParamFromPath(g, "id")
	if err != nil {
		api_helper.ClientErrorHandler(g, 40004)
		return
	}
	current, err := c.ipService.FindIPByID(id)
	if err != nil {
		api_helper.ClientErrorHandler(g, 40005)
		return
	}
	err = c.ipService.DeleteIP(current)
	if err != nil {
		api_helper.InternalErrorHandler(g, err)
		return
	}
	g.JSON(http.StatusOK, api_helper.Response{Status: 0, Msg: "操作成功！"})
}
