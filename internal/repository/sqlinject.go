package repository

import (
	"GoWAFer/internal/model"
	"fmt"
	"gorm.io/gorm"
)

type SqlInjectRepository struct {
	db *gorm.DB
}

func NewSqlInjectRepository(db *gorm.DB) *SqlInjectRepository {
	return &SqlInjectRepository{db: db}
}

// Create 新增sql注入规则
func (r *SqlInjectRepository) Create(l *model.SqlInjectionRules) error {
	return r.db.Create(l).Error
}

// FindPaginated 分页查询sql注入规则
func (r *SqlInjectRepository) FindPaginated(pageIndex, pageSize int, keyword string) ([]model.SqlInjectionRules, int) {
	var rules []model.SqlInjectionRules
	var count int64
	query := r.db.Model(&model.SqlInjectionRules{})

	if keyword != "" {
		query = query.Where("Rule LIKE ?", fmt.Sprintf("%%%s%%", keyword))
	}

	query.Count(&count).Limit(pageSize).Offset((pageIndex - 1) * pageSize).Find(&rules)

	return rules, int(count)
}

func (r *SqlInjectRepository) FindByID(id uint) (*model.SqlInjectionRules, error) {
	var rule *model.SqlInjectionRules
	if err := r.db.First(&rule, id).Error; err != nil {
		return nil, err
	}
	return rule, nil
}

// Update 编辑sql注入规则
func (r *SqlInjectRepository) Update(l *model.SqlInjectionRules) error {
	return r.db.Save(l).Error
}

// Delete 删除sql注入规则
func (r *SqlInjectRepository) Delete(l *model.SqlInjectionRules) error {
	return r.db.Delete(l).Error
}

func (r *SqlInjectRepository) FindAll() []model.SqlInjectionRules {
	var rules []model.SqlInjectionRules
	r.db.Find(&rules)
	return rules
}
