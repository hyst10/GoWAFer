package service

import (
	"GoWAFer/internal/model"
	"GoWAFer/internal/repository"
)

type UserService struct {
	userRepository *repository.UserRepository
}

func NewUserService(r *repository.UserRepository) *UserService {
	r.Migrate()
	r.InsertDefaultUser()
	return &UserService{
		userRepository: r,
	}
}

func (c *UserService) FindUserByUsername(username string) (*model.User, error) {
	return c.userRepository.FindByUsername(username)
}

func (c *UserService) UpdateUserInfo(u *model.User) error {
	return c.userRepository.Update(u)
}
