package service

import (
	"GoWAFer/internal/model"
	"GoWAFer/internal/repository"
	"GoWAFer/pkg/pagination"
)

type IPService struct {
	ipRepository *repository.IPRepository
}

func NewIPListService(r *repository.IPRepository) *IPService {
	return &IPService{ipRepository: r}
}

// CreateIP 创建IP
func (c *IPService) CreateIP(i *model.IP) error {
	return c.ipRepository.Create(i)
}

// FindPaginatedIPs 分页查询IP
func (c *IPService) FindPaginatedIPs(page *pagination.Pages, ipType, keyword string) *pagination.Pages {
	items, count := c.ipRepository.FindPaginated(page.Page, page.PerPage, ipType, keyword)
	page.Items = items
	page.Total = count
	return page
}

func (c *IPService) FindIPByID(id uint) (*model.IP, error) {
	return c.ipRepository.FindByID(id)
}

// UpdateIP 编辑IP
func (c *IPService) UpdateIP(i *model.IP) error {
	return c.ipRepository.Update(i)
}

// DeleteIP 删除IP
func (c *IPService) DeleteIP(i *model.IP) error {
	return c.ipRepository.Delete(i)
}
