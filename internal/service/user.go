package service

import (
	"golang/stockLkBack/internal/model"
	"golang/stockLkBack/internal/repository"
)

type UsersService struct {
	repo repository.User
}

func NewUsersService(repo repository.User) *UsersService {
	return &UsersService{repo: repo}
}

func (s *UsersService) Create(user model.UserCreateBody) (*model.User, error) {
	return s.repo.Create(user)
}

func (s *UsersService) GetAll() ([]model.User, error) {
	return s.repo.GetAll()
}

func (s *UsersService) GetByID(id int) (*model.User, error) {
	return s.repo.GetByID(id)
}

func (s *UsersService) Delete(id int) error {
	return s.repo.Delete(id)
}

func (s *UsersService) Update(id int, user model.UserEditBody) (*model.User, error) {
	return s.repo.Update(id, user)
}

func (s *UsersService) Login(user model.LoginRequest) (*model.TokenSuccess, error) {
	return s.repo.Login(user)
}

func (s *UsersService) ChangeUserRole(id int, userRoleReq model.UserRoleBody) (*model.User, error) {
	return s.repo.ChangeUserRole(id, userRoleReq)
}

func (s *UsersService) ChangePassword(id int, changePassworReq model.UserChangePasswordBody) (*model.Success, error) {
	return s.repo.ChangePassword(id, changePassworReq)
}
