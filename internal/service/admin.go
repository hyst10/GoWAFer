package service

import (
	"GoWAFer/internal/model"
	"GoWAFer/internal/repository"
)

// AdminService 管理员Service层接口
type AdminService struct {
	adminRepository *repository.AdminRepository
}

// NewAdminService 实例化管理员Service层接口
func NewAdminService(r *repository.AdminRepository) *AdminService {
	return &AdminService{
		adminRepository: r,
	}
}

// FindAdminByUsername 通过用户名查询管理员
func (c *AdminService) FindAdminByUsername(username string) (*model.Admin, error) {
	return c.adminRepository.FindByUsername(username)
}

// UpdateAdminInfo 更新管理员信息
func (c *AdminService) UpdateAdminInfo(u *model.Admin) error {
	return c.adminRepository.Update(u)
}
