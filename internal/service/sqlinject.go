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

func (c *SqlInjectService) CreateRule(r *model.SqlInjectionRules) error {
	return c.sqlInjectRepository.Create(r)
}

func (c *SqlInjectService) FindPaginatedRules(page *pagination.Pages, keyword string) *pagination.Pages {
	items, count := c.sqlInjectRepository.FindPaginated(page.Page, page.PerPage, keyword)
	page.Items = items
	page.Total = count
	return page
}

func (c *SqlInjectService) FindRuleByID(id uint) (*model.SqlInjectionRules, error) {
	return c.sqlInjectRepository.FindByID(id)
}

func (c *SqlInjectService) UpdateRule(r *model.SqlInjectionRules) error {
	return c.sqlInjectRepository.Update(r)
}

func (c *SqlInjectService) DeleteRule(r *model.SqlInjectionRules) error {
	return c.sqlInjectRepository.Delete(r)
}
