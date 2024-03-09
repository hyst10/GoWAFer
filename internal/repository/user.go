package repository

import (
	"GoWAFer/internal/model"
	"gorm.io/gorm"
)

// UserRepository 用户仓库接口
type UserRepository struct {
	db *gorm.DB
}

// NewUserRepository 实例化用户仓库接口
func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) FindByUsername(username string) (*model.User, error) {
	var user model.User
	err := r.db.Where("Username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, err
}

func (r *UserRepository) Update(u *model.User) error {
	return r.db.Save(u).Error
}
