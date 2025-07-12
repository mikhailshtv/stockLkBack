package service

import (
	"context"
	"errors"
	"log"

	"github.com/mikhailshtv/stockLkBack/internal/model"
	"github.com/mikhailshtv/stockLkBack/internal/repository"
)

const (
	logUsersTableName = "logUser"
)

type UsersService struct {
	repo repository.User
	ctx  context.Context
}

func NewUsersService(ctx context.Context, repo repository.User) *UsersService {
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
	createdUser, err := s.repo.Create(s.ctx, user)
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
	return s.repo.GetByID(s.ctx, id)
}

func (s *UsersService) Delete(id int) error {
	delitedUser, err := s.repo.Delete(s.ctx, id)
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
	updatedUser, err := s.repo.Update(s.ctx, id, user)
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
	return s.repo.Login(s.ctx, user)
}

func (s *UsersService) ChangeUserRole(id int, userRoleReq model.UserRoleBody) (*model.User, error) {
	return s.repo.ChangeUserRole(s.ctx, id, userRoleReq)
}

func (s *UsersService) ChangePassword(id int, changePassworReq model.UserChangePasswordBody) (*model.Success, error) {
	return s.repo.ChangePassword(s.ctx, id, changePassworReq)
}
