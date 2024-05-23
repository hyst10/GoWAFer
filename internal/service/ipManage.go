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
func (s *IPManageService) AddIP(ip string, expiration int, isBlack bool) error {
	return s.ipManageRepository.Add(ip, expiration, isBlack)
}

// DeleteIP 删除一条IP记录
func (s *IPManageService) DeleteIP(ip string, isBlack bool) error {
	return s.ipManageRepository.Del(ip, isBlack)
}

// GetIPWithPagination 分页查询路由管理记录
func (s *IPManageService) GetIPWithPagination(page *pagination.Pages, isBlack bool, query string) *pagination.Pages {
	items, count := s.ipManageRepository.GetAllWithPagination(page.Page, page.PerPage, isBlack, query)
	page.Items = items
	page.Total = count
	return page
}
