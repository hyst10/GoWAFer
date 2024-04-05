package service

import (
	"GoWAFer/internal/repository"
	"GoWAFer/internal/types"
)

type XssDetectService struct {
	xssDetectRepository *repository.XssDetectRepository
}

func NewXssDetectService(r *repository.XssDetectRepository) *XssDetectService {
	return &XssDetectService{xssDetectRepository: r}
}

// AddRule 新增xss防护规则
func (c *XssDetectService) AddRule(rule string) error {
	return c.xssDetectRepository.Add(rule)
}

// GetAllRules 获取全部防护规则
func (c *XssDetectService) GetAllRules() ([]types.SqlInjectRule, int) {
	return c.xssDetectRepository.GetAll()
}

// DeleteRule 删除防护规则
func (c *XssDetectService) DeleteRule(rule string) error {
	return c.xssDetectRepository.Delete(rule)
}
