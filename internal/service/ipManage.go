package service

import (
	"GoWAFer/internal/repository"
	"GoWAFer/pkg/pagination"
)

// IPManageService IP管理服务
type IPManageService struct {
	ipManageRepository *repository.IPManageRepository
}

// NewIPManageService 实例化IP管理服务
func NewIPManageService(r *repository.IPManageRepository) *IPManageService {
	return &IPManageService{ipManageRepository: r}
}

// AddIP 添加一条IP记录
func (c *IPManageService) AddIP(ip string, expiration, ipType int) error {
	return c.ipManageRepository.Add(ip, expiration, ipType)
}

// DeleteIP 删除一条IP记录
func (c *IPManageService) DeleteIP(ip string, ipType int) error {
	return c.ipManageRepository.Del(ip, ipType)
}

// GetIPWithPagination 分页查询路由管理记录
func (c *IPManageService) GetIPWithPagination(page *pagination.Pages, ipType int, keyword string) *pagination.Pages {
	items, count := c.ipManageRepository.GetAllWithPagination(page.Page, page.PerPage, ipType, keyword)
	page.Items = items
	page.Total = count
	return page
}
