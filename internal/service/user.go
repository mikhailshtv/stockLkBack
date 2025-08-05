package service

import (
	"context"
	"github.com/mikhailshtv/stockLkBack/pkg/errors"

	"github.com/mikhailshtv/stockLkBack/internal/model"
	"github.com/mikhailshtv/stockLkBack/internal/repository"
	"github.com/mikhailshtv/stockLkBack/pkg/logger"

	"go.uber.org/zap"
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
			logger.GetLogger().Error("failed to hash password",
				zap.Error(err),
				zap.String("email", userRequest.Email),
			)
			return nil, errors.NewValidationError("ошибка хеширования пароля", err)
		}
	} else {
		return nil, errors.NewValidationError("ошибка подтверждения пароля", nil)
	}
	createdUser, err := s.repo.Create(s.ctx, user)
	var result any
	var status string
	if err != nil {
		logger.GetLogger().Error("failed to create user in repository",
			zap.Error(err),
			zap.String("email", userRequest.Email),
		)
		result = model.Error{Error: err.Error()}
		status = logErrorStatus
	} else {
		logger.GetLogger().Info("user created successfully",
			zap.Int("user_id", createdUser.ID),
			zap.String("email", createdUser.Email),
		)
		result = createdUser
		status = logSuccessStatus
	}

	_, logErr := s.repo.WriteLog(result, "Create", status, logUsersTableName)
	if logErr != nil {
		logger.GetLogger().Error("failed to write log for user creation",
			zap.Error(logErr),
		)
	}
	return createdUser, err
}

func (s *UsersService) GetAll() ([]model.User, error) {
	users, err := s.repo.GetAll(s.ctx)
	if err != nil {
		logger.GetLogger().Error("failed to get users from repository",
			zap.Error(err),
		)
		return nil, errors.NewDatabaseError("ошибка получения списка пользователей", err)
	}
	return users, nil
}

func (s *UsersService) GetByID(id int) (*model.User, error) {
	user, err := s.repo.GetByID(s.ctx, id)
	if err != nil {
		logger.GetLogger().Error("failed to get user by ID from repository",
			zap.Error(err),
			zap.Int("user_id", id),
		)
		if err.Error() == "пользователь не найден" {
			return nil, errors.NewNotFoundError("пользователь", err)
		}
		return nil, errors.NewDatabaseError("ошибка получения пользователя", err)
	}
	return user, nil
}

func (s *UsersService) Delete(id int) error {
	deletedUser, err := s.repo.Delete(s.ctx, id)
	var result any
	var status string
	if err != nil {
		logger.GetLogger().Error("failed to delete user from repository",
			zap.Error(err),
			zap.Int("user_id", id),
		)
		result = model.Error{Error: err.Error()}
		status = logErrorStatus
	} else {
		logger.GetLogger().Info("user deleted successfully",
			zap.Int("user_id", id),
		)
		result = deletedUser
		status = logSuccessStatus
	}

	_, logErr := s.repo.WriteLog(result, "Delete", status, logUsersTableName)
	if logErr != nil {
		logger.GetLogger().Error("failed to write log for user deletion",
			zap.Error(logErr),
		)
	}
	return err
}

func (s *UsersService) Update(id int, user model.UserEditBody) (*model.User, error) {
	updatedUser, err := s.repo.Update(s.ctx, id, user)
	var result any
	var status string
	if err != nil {
		logger.GetLogger().Error("failed to update user in repository",
			zap.Error(err),
			zap.Int("user_id", id),
		)
		result = model.Error{Error: err.Error()}
		status = logErrorStatus
	} else {
		logger.GetLogger().Info("user updated successfully",
			zap.Int("user_id", id),
			zap.String("email", updatedUser.Email),
		)
		result = updatedUser
		status = logSuccessStatus
	}

	_, logErr := s.repo.WriteLog(result, "Update", status, logUsersTableName)
	if logErr != nil {
		logger.GetLogger().Error("failed to write log for user update",
			zap.Error(logErr),
		)
	}
	return updatedUser, err
}

func (s *UsersService) Login(user model.LoginRequest) (*model.TokenSuccess, error) {
	token, err := s.repo.Login(s.ctx, user)
	if err != nil {
		logger.GetLogger().Error("failed to login user",
			zap.Error(err),
			zap.String("login", user.Login),
		)
		return nil, err
	}
	logger.GetLogger().Info("user logged in successfully",
		zap.String("login", user.Login),
	)
	return token, nil
}

func (s *UsersService) ChangeUserRole(id int, userRoleReq model.UserRoleBody) (*model.User, error) {
	updatedUser, err := s.repo.ChangeUserRole(s.ctx, id, userRoleReq)
	if err != nil {
		logger.GetLogger().Error("failed to change user role",
			zap.Error(err),
			zap.Int("user_id", id),
		)
		return nil, err
	}
	logger.GetLogger().Info("user role changed successfully",
		zap.Int("user_id", id),
		zap.String("new_role", string(userRoleReq.Role)),
	)
	return updatedUser, nil
}

func (s *UsersService) ChangePassword(id int, changePasswordReq model.UserChangePasswordBody) (*model.Success, error) {
	result, err := s.repo.ChangePassword(s.ctx, id, changePasswordReq)
	if err != nil {
		logger.GetLogger().Error("failed to change user password",
			zap.Error(err),
			zap.Int("user_id", id),
		)
		return nil, err
	}
	logger.GetLogger().Info("user password changed successfully",
		zap.Int("user_id", id),
	)
	return result, nil
}
