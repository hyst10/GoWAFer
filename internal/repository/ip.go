package repository

import (
	"GoWAFer/internal/model"
	"fmt"
	"gorm.io/gorm"
)

type IPRepository struct {
	db *gorm.DB
}

func NewIPRepository(db *gorm.DB) *IPRepository {
	return &IPRepository{db: db}
}

// FindPaginated 分页查询IP
func (r *IPRepository) FindPaginated(pageIndex, pageSize int, ipType, keyword string) ([]model.IP, int) {
	var ips []model.IP
	var count int64
	query := r.db.Model(&model.IP{}).Where("Type = ?", ipType)

	if keyword != "" {
		query = query.Where("IPAddress LIKE ?", fmt.Sprintf("%%%s%%", keyword))
	}

	query.Count(&count).Limit(pageSize).Offset((pageIndex - 1) * pageSize).Find(&ips)

	return ips, int(count)
}

// FindByID 通过ID查询IP
func (r *IPRepository) FindByID(id uint) (*model.IP, error) {
	var ip *model.IP
	err := r.db.First(&ip, id).Error
	if err != nil {
		return nil, err
	}
	return ip, err
}

// Create 创建IP
func (r *IPRepository) Create(i *model.IP) error {
	return r.db.Create(i).Error
}

// Update 编辑IP
func (r *IPRepository) Update(i *model.IP) error {
	return r.db.Save(i).Error
}

// Delete 删除IP
func (r *IPRepository) Delete(i *model.IP) error {
	return r.db.Delete(i).Error
}
