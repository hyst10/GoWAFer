package service

import (
	"GoWAFer/internal/repository"
	"GoWAFer/internal/types"
)

type SqlInjectService struct {
	sqlInjectRepository *repository.SqlInjectRepository
}

func NewSqlInjectService(r *repository.SqlInjectRepository) *SqlInjectService {
	return &SqlInjectService{sqlInjectRepository: r}
}

// AddRule 新增sql注入防护规则
func (c *SqlInjectService) AddRule(rule string) error {
	return c.sqlInjectRepository.Add(rule)
}

// GetAllRules 获取全部防护规则
func (c *SqlInjectService) GetAllRules() ([]types.SqlInjectRule, int) {
	return c.sqlInjectRepository.GetAll()
}

// DeleteRule 删除防护规则
func (c *SqlInjectService) DeleteRule(rule string) error {
	return c.sqlInjectRepository.Delete(rule)
}
