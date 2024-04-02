package service

import (
	"GoWAFer/internal/repository"
	"GoWAFer/pkg/pagination"
)

// RoutingManageService 路由管理服务
type RoutingManageService struct {
	routingManageRepository *repository.RoutingManageRepository
}

// NewRoutingManageService 实例化路由管理服务
func NewRoutingManageService(r *repository.RoutingManageRepository) *RoutingManageService {
	return &RoutingManageService{routingManageRepository: r}
}

// AddRouting 添加一条路由记录
func (c *RoutingManageService) AddRouting(routing, method string, routingType int) error {
	return c.routingManageRepository.Add(routing, method, routingType)
}

// DeleteRouting 删除一条路由记录
func (c *RoutingManageService) DeleteRouting(routing, method string, routingType int) error {
	return c.routingManageRepository.Del(routing, method, routingType)
}

// GetRoutingWithPagination 分页查询路由管理记录
func (c *RoutingManageService) GetRoutingWithPagination(page *pagination.Pages, routerType int, keyword string) *pagination.Pages {
	items, count := c.routingManageRepository.GetAllWithPagination(page.Page, page.PerPage, routerType, keyword)
	page.Items = items
	page.Total = count
	return page
}
