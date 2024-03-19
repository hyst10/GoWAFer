package service

import (
	"GoWAFer/internal/model"
	"GoWAFer/internal/repository"
	"GoWAFer/pkg/pagination"
)

type XssDetectService struct {
	xssDetectRepository *repository.XssDetectRepository
}

func NewXssDetectService(r *repository.XssDetectRepository) *XssDetectService {
	return &XssDetectService{xssDetectRepository: r}
}

// CreateRule 创建xss防护规则
func (c *XssDetectService) CreateRule(r *model.XssDetectRules) error {
	return c.xssDetectRepository.Create(r)
}

// FindPaginatedRules 分页查询xss防护规则
func (c *XssDetectService) FindPaginatedRules(page *pagination.Pages, keyword string) *pagination.Pages {
	items, count := c.xssDetectRepository.FindPaginated(page.Page, page.PerPage, keyword)
	page.Items = items
	page.Total = count
	return page
}

// FindRuleByID 通过主键ID查询防护规则
func (c *XssDetectService) FindRuleByID(id uint) (*model.XssDetectRules, error) {
	return c.xssDetectRepository.FindByID(id)
}

// UpdateRule 更新防护规则
func (c *XssDetectService) UpdateRule(r *model.XssDetectRules) error {
	return c.xssDetectRepository.Update(r)
}

// DeleteRule 删除防护规则
func (c *XssDetectService) DeleteRule(r *model.XssDetectRules) error {
	return c.xssDetectRepository.Delete(r)
}
