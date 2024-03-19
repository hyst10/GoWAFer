package repository

import (
	"GoWAFer/internal/model"
	"gorm.io/gorm"
)

// AdminRepository 管理员Repository层接口
type AdminRepository struct {
	db *gorm.DB
}

// NewAdminRepository 实例化Repository层接口
func NewAdminRepository(db *gorm.DB) *AdminRepository {
	return &AdminRepository{db: db}
}

// FindByUsername 通过用户名查找管理员
func (r *AdminRepository) FindByUsername(username string) (*model.Admin, error) {
	var admin model.Admin
	err := r.db.Where("Username = ?", username).First(&admin).Error
	if err != nil {
		return nil, err
	}
	return &admin, err
}

// FindByID 通过主键ID查找管理员
func (r *AdminRepository) FindByID(id uint) (*model.Admin, error) {
	var admin model.Admin
	if err := r.db.First(&admin, id).Error; err != nil {
		return nil, err
	}
	return &admin, nil
}

// Update 更新管理员信息
func (r *AdminRepository) Update(u *model.Admin) error {
	return r.db.Save(u).Error
}
