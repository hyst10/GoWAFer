package repository

import (
	"GoWAFer/internal/model"
	"GoWAFer/pkg/hash_handler"
	"fmt"
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

// Migrate 迁移管理员表
func (r *UserRepository) Migrate() {
	err := r.db.AutoMigrate(&model.User{})
	if err != nil {
		panic(fmt.Sprintf("管理员表创建失败：%v", err))
	}
}

// InsertDefaultUser 检查是否管理员表中是否存在用户，不存在则创建一个默认管理员
func (r *UserRepository) InsertDefaultUser() {
	// 哈希加盐加密
	defaultPassword, _ := hash_handler.EncryptPassword("123456")
	defaultUser := model.User{Model: gorm.Model{ID: 1}, Username: "admin", Password: defaultPassword, Nickname: "admin"}
	err := r.db.FirstOrCreate(&defaultUser, model.User{Model: gorm.Model{ID: 1}}).Error
	if err != nil {
		panic(fmt.Sprintf("默认管理员创建失败：%v", err))
	}
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
