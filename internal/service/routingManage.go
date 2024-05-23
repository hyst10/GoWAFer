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
func (c *RoutingManageService) AddRouting(path string, isBlack bool) error {
	return c.routingManageRepository.Add(path, isBlack)
}

// DeleteRouting 删除一条路由记录
func (c *RoutingManageService) DeleteRouting(path string, isBlack bool) error {
	return c.routingManageRepository.Del(path, isBlack)
}

// GetRoutingWithPagination 分页查询路由管理记录
func (c *RoutingManageService) GetRoutingWithPagination(page *pagination.Pages, isBlack bool, query string) *pagination.Pages {
	items, count := c.routingManageRepository.GetAllWithPagination(page.Page, page.PerPage, isBlack, query)
	page.Items = items
	page.Total = count
	return page
}
