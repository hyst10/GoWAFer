package repository

import (
	"GoWAFer/internal/model"
	"fmt"
	"gorm.io/gorm"
)

type XssDetectRepository struct {
	db *gorm.DB
}

func NewXssDetectRepository(db *gorm.DB) *XssDetectRepository {
	return &XssDetectRepository{db: db}
}

// Create 新增xss防护规则
func (r *XssDetectRepository) Create(l *model.XssDetectRules) error {
	return r.db.Create(l).Error
}

// FindPaginated 分页查询xss防护规则
func (r *XssDetectRepository) FindPaginated(pageIndex, pageSize int, keyword string) ([]model.XssDetectRules, int) {
	var rules []model.XssDetectRules
	var count int64
	query := r.db.Model(&model.XssDetectRules{})

	if keyword != "" {
		query = query.Where("Rule LIKE ?", fmt.Sprintf("%%%s%%", keyword))
	}

	query.Count(&count).Limit(pageSize).Offset((pageIndex - 1) * pageSize).Find(&rules)

	return rules, int(count)
}

func (r *XssDetectRepository) FindByID(id uint) (*model.XssDetectRules, error) {
	var rule *model.XssDetectRules
	if err := r.db.First(&rule, id).Error; err != nil {
		return nil, err
	}
	return rule, nil
}

// Update 编辑xss防护规则
func (r *XssDetectRepository) Update(l *model.XssDetectRules) error {
	return r.db.Save(l).Error
}

// Delete 删除xss防护规则
func (r *XssDetectRepository) Delete(l *model.XssDetectRules) error {
	return r.db.Delete(l).Error
}

func (r *XssDetectRepository) FindAll() []model.XssDetectRules {
	var rules []model.XssDetectRules
	r.db.Find(&rules)
	return rules
}
