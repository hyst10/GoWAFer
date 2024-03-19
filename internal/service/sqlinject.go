package service

import (
	"GoWAFer/internal/model"
	"GoWAFer/internal/repository"
	"GoWAFer/pkg/pagination"
)

type SqlInjectService struct {
	sqlInjectRepository *repository.SqlInjectRepository
}

func NewSqlInjectService(r *repository.SqlInjectRepository) *SqlInjectService {
	return &SqlInjectService{sqlInjectRepository: r}
}

// CreateRule 创建sql注入防护规则
func (c *SqlInjectService) CreateRule(r *model.SqlInjectionRules) error {
	return c.sqlInjectRepository.Create(r)
}

// FindPaginatedRules 分页查询sql注入防护规则
func (c *SqlInjectService) FindPaginatedRules(page *pagination.Pages, keyword string) *pagination.Pages {
	items, count := c.sqlInjectRepository.FindPaginated(page.Page, page.PerPage, keyword)
	page.Items = items
	page.Total = count
	return page
}

// FindRuleByID 通过主键ID查询防护规则
func (c *SqlInjectService) FindRuleByID(id uint) (*model.SqlInjectionRules, error) {
	return c.sqlInjectRepository.FindByID(id)
}

// UpdateRule 更新防护规则
func (c *SqlInjectService) UpdateRule(r *model.SqlInjectionRules) error {
	return c.sqlInjectRepository.Update(r)
}

// DeleteRule 删除防护规则
func (c *SqlInjectService) DeleteRule(r *model.SqlInjectionRules) error {
	return c.sqlInjectRepository.Delete(r)
}
