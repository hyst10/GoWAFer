package service

import (
	"GoWAFer/internal/model"
	"GoWAFer/internal/repository"
	"GoWAFer/pkg/pagination"
)

type RoutingService struct {
	routingRepository *repository.RoutingRepository
}

func NewRoutingService(r *repository.RoutingRepository) *RoutingService {
	return &RoutingService{routingRepository: r}
}

// CreateRouting 添加路由名单
func (c *RoutingService) CreateRouting(r *model.Routing) error {
	return c.routingRepository.Create(r)
}

// FindPaginatedRouters 分页查询路由
func (c *RoutingService) FindPaginatedRouters(page *pagination.Pages, routerType, keyword string) *pagination.Pages {
	items, count := c.routingRepository.FindPaginated(page.Page, page.PerPage, routerType, keyword)
	page.Items = items
	page.Total = count
	return page
}

// FindRoutingByID 通过主键ID查找路由
func (c *RoutingService) FindRoutingByID(id uint) (*model.Routing, error) {
	return c.routingRepository.FindByID(id)
}

// UpdateRouting 编辑路由
func (c *RoutingService) UpdateRouting(r *model.Routing) error {
	return c.routingRepository.Update(r)
}

// DeleteRouting 删除路由
func (c *RoutingService) DeleteRouting(r *model.Routing) error {
	return c.routingRepository.Delete(r)
}
