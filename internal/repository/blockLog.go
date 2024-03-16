package repository

import (
	"GoWAFer/internal/model"
	"fmt"
	"gorm.io/gorm"
)

type BlockLogRepository struct {
	db *gorm.DB
}

func NewBlockLogRepository(db *gorm.DB) *BlockLogRepository {
	return &BlockLogRepository{db: db}
}

// Create 新增拦截日志记录
func (r *BlockLogRepository) Create(l *model.BlockLog) error {
	return r.db.Create(l).Error
}

// FindPaginated 分页查询
func (r *BlockLogRepository) FindPaginated(pageIndex, pageSize int, keyword string) ([]model.BlockLog, int) {
	var logs []model.BlockLog
	var count int64
	query := r.db.Model(&model.BlockLog{}).Order("CreatedAt DESC")

	if keyword != "" {
		query = query.Where("ClientIP LIKE ?", fmt.Sprintf("%%%s%%", keyword))
	}

	query.Count(&count).Limit(pageSize).Offset((pageIndex - 1) * pageSize).Find(&logs)

	return logs, int(count)
}
