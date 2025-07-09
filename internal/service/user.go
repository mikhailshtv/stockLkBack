package service

import (
	"context"
	"errors"
	"golang/stockLkBack/internal/model"
	"golang/stockLkBack/internal/repository"
	"log"
)

const (
	logUsersTableName = "logUser"
)

type UsersService struct {
	repo repository.User
	ctx  context.Context
}

func NewUsersService(repo repository.User, ctx context.Context) *UsersService {
	return &UsersService{repo: repo, ctx: ctx}
}

func (s *UsersService) Create(userRequest model.UserCreateBody) (*model.User, error) {
	var user model.User
	user.FirstName = userRequest.FirstName
	user.LastName = userRequest.LastName
	user.Login = userRequest.Login
	user.Email = userRequest.Email
	if userRequest.Password == userRequest.PasswordConfirm {
		err := user.HashPassword(userRequest.Password)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, errors.New("ошибка подтверждения пароля")
	}
	createdUser, err := s.repo.Create(user, s.ctx)
	var result any
	var status string
	if err != nil {
		result = model.Error{Error: err.Error()}
		status = logErrorStatus
	} else {
		result = createdUser
		status = logSuccessStatus
	}

	_, logErr := s.repo.WriteLog(result, "Create", status, logUsersTableName)
	if logErr != nil {
		log.Println(logErr.Error())
	}
	return createdUser, err
}

func (s *UsersService) GetAll() ([]model.User, error) {
	return s.repo.GetAll(s.ctx)
}

func (s *UsersService) GetByID(id int) (*model.User, error) {
	return s.repo.GetByID(id, s.ctx)
}

func (s *UsersService) Delete(id int) error {
	delitedUser, err := s.repo.Delete(id, s.ctx)
	var result any
	var status string
	if err != nil {
		result = model.Error{Error: err.Error()}
		status = logErrorStatus
	} else {
		result = delitedUser
		status = logSuccessStatus
	}

	_, logErr := s.repo.WriteLog(result, "Delete", status, logUsersTableName)
	if logErr != nil {
		log.Println(logErr.Error())
	}
	return err
}

func (s *UsersService) Update(id int, user model.UserEditBody) (*model.User, error) {
	updatedUser, err := s.repo.Update(id, user, s.ctx)
	var result any
	var status string
	if err != nil {
		result = model.Error{Error: err.Error()}
		status = logErrorStatus
	} else {
		result = updatedUser
		status = logSuccessStatus
	}

	_, logErr := s.repo.WriteLog(result, "Update", status, logUsersTableName)
	if logErr != nil {
		log.Println(logErr.Error())
	}
	return updatedUser, err
}

func (s *UsersService) Login(user model.LoginRequest) (*model.TokenSuccess, error) {
	return s.repo.Login(user, s.ctx)
}

func (s *UsersService) ChangeUserRole(id int, userRoleReq model.UserRoleBody) (*model.User, error) {
	return s.repo.ChangeUserRole(id, userRoleReq, s.ctx)
}

func (s *UsersService) ChangePassword(id int, changePassworReq model.UserChangePasswordBody) (*model.Success, error) {
	return s.repo.ChangePassword(id, changePassworReq, s.ctx)
}
