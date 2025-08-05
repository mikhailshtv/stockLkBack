package handler

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/mikhailshtv/stockLkBack/internal/middleware"
	"github.com/mikhailshtv/stockLkBack/internal/model"
	"github.com/mikhailshtv/stockLkBack/pkg/errors"
	"github.com/mikhailshtv/stockLkBack/pkg/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// CreateUser
// @Summary Создание/регистрация пользователя
// @Tags Users
// @Accept			json
// @Produce		json
// @Param user body model.UserCreateBody true "Объект пользователя"
// @Success 201 {object} model.User "Created"
// @Failure 400 {object} model.Error "Invalid request"
// @Failure 500 {object} model.Error "Internal"
// @Router /api/v1/users [post].
func (h *Handler) CreateUser(ctx *gin.Context) {
	var userReq model.UserCreateBody
	if err := ctx.ShouldBindJSON(&userReq); err != nil {
		middleware.HandleError(ctx, errors.NewValidationError("Неверный формат данных", err))
		return
	}
	user, err := h.Services.User.Create(userReq)
	if err != nil {
		logger.GetLogger().Error("ошибка создания пользователя",
			zap.Error(err),
			zap.String("email", userReq.Email),
		)
		middleware.HandleError(ctx, err)
		return
	}
	ctx.JSON(http.StatusCreated, user)
}

// Login
// @Summary Аутентификация пользователя
// @Tags Login
// @Accept			json
// @Produce		json
// @Param user body model.LoginRequest true "Данные для аутентификации пользователя"
// @Success 200 {object} model.TokenSuccess
// @Failure 400 {object} model.Error "Invalid request"
// @Failure 401 {object} model.Error "Unauthorized"
// @Failure 500 {object} model.Error "Internal server error"
// @Router /api/v1/login [post].
func (h *Handler) Login(ctx *gin.Context) {
	var loginRequest model.LoginRequest
	if err := ctx.ShouldBindJSON(&loginRequest); err != nil {
		middleware.HandleError(ctx, errors.NewValidationError("Неверный формат данных", err))
		return
	}
	TokenSuccess, err := h.Services.User.Login(loginRequest)
	if err != nil {
		var appErr *errors.AppError
		switch err.Error() {
		case "логин или пароль пользователя недействителен":
			appErr = errors.NewUnauthorizedError("Логин или пароль пользователя недействителен", err)
		case "ошибка генерации токена":
			appErr = errors.NewInternalError("Ошибка генерации токена", err)
		default:
			appErr = errors.NewInternalError("Ошибка аутентификации", err)
		}

		logger.GetLogger().Error("login failed",
			zap.Error(err),
			zap.String("login", loginRequest.Login),
		)
		middleware.HandleError(ctx, appErr)
		return
	}
	ctx.JSON(http.StatusOK, TokenSuccess)
}

// EditUser
// @Summary Редактирование пользователя
// @Tags Users
// @Accept			json
// @Produce		json
// @Param user body model.UserEditBody true "Объект пользователя"
// @Success 200 {object} model.User
// @Failure 400 {object} model.Error "Invalid request"
// @Failure 401 {object} model.Error "Unauthorized"
// @Failure 404 {object} model.Error "Not found"
// @Failure 500 {object} model.Error "Internal"
// @Param id path string true "id пользователя"
// @Router /api/v1/users/{id} [put]
// @Security BearerAuth.
func (h *Handler) EditUser(ctx *gin.Context) {
	idStr := ctx.Params.ByName("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		middleware.HandleError(ctx, errors.NewValidationError("Некорректный ID пользователя", err))
		return
	}
	var userEdit model.UserEditBody
	if err := ctx.ShouldBindJSON(&userEdit); err != nil {
		middleware.HandleError(ctx, errors.NewValidationError("Неверный формат данных", err))
		return
	}
	if userEdit.Email == "" && userEdit.FirstName == "" && userEdit.LastName == "" {
		middleware.HandleError(ctx, errors.NewValidationError(
			"Тело запроса должно содержать хотя бы одно поле для обновления",
			nil,
		))
		return
	}
	user, err := h.Services.User.Update(id, userEdit)
	if err != nil {
		logger.GetLogger().Error("failed to edit user",
			zap.Error(err),
			zap.Int("user_id", id),
		)
		if strings.Contains(err.Error(), "пользователь не найден") {
			middleware.HandleError(ctx, errors.NewNotFoundError("пользователь", err))
			return
		}
		middleware.HandleError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, user)
}

// ChangeUserRole
// @Summary Изменение роли пользователя
// @Tags Users
// @Accept			json
// @Produce		json
// @Param user body model.UserRoleBody true "Объект с ролью пользователя"
// @Success 200 {object} model.User
// @Failure 400 {object} model.Error "Invalid request"
// @Failure 401 {object} model.Error "Unauthorized"
// @Failure 404 {object} model.Error "Not found"
// @Failure 500 {object} model.Error "Internal"
// @Param id path string true "id пользователя"
// @Router /api/v1/users/{id}/role [patch]
// @Security BearerAuth.
func (h *Handler) ChangeUserRole(ctx *gin.Context) {
	role, exists := ctx.Get(userRoleKey)
	if !exists {
		middleware.HandleError(ctx, errors.NewValidationError("Некорректная роль пользователя", nil))
		return
	}
	isEmployee := role == model.RoleEmployee
	if !isEmployee {
		middleware.HandleError(ctx, errors.NewForbiddenError("Недостаточно прав для выполнения операции", nil))
		return
	}
	idStr := ctx.Params.ByName("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		middleware.HandleError(ctx, errors.NewValidationError("Некорректный ID пользователя", err))
		return
	}
	var userRole *model.UserRoleBody
	if err := ctx.ShouldBindJSON(&userRole); err != nil {
		middleware.HandleError(ctx, errors.NewValidationError("Неверный формат данных", err))
		return
	}
	if !userRole.Role.Valid() {
		middleware.HandleError(ctx, errors.NewValidationError("Роль не существует", err))
		return
	}

	user, err := h.Services.User.ChangeUserRole(id, *userRole)
	if err != nil {
		logger.GetLogger().Error("failed to change user role",
			zap.Error(err),
			zap.Int("user_id", id),
		)
		if strings.Contains(err.Error(), "пользователь не найден") {
			middleware.HandleError(ctx, errors.NewNotFoundError("пользователь", err))
			return
		}
		middleware.HandleError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, user)
}

// ChangeUserPassword
// @Summary Изменение пароля пользователя
// @Tags Users
// @Accept			json
// @Produce		json
// @Param user body model.UserChangePasswordBody true "Объект с паролем пользователя"
// @Success 200 {object} model.Success
// @Failure 400 {object} model.Error "Invalid request"
// @Failure 401 {object} model.Error "Unauthorized"
// @Failure 404 {object} model.Error "Not found"
// @Failure 500 {object} model.Error "Internal"
// @Param id path string true "id пользователя"
// @Router /api/v1/users/{id}/password [patch]
// @Security BearerAuth.
func (h *Handler) ChangeUserPassword(ctx *gin.Context) {
	idStr := ctx.Params.ByName("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		middleware.HandleError(ctx, errors.NewValidationError("Некорректный ID пользователя", err))
		return
	}
	var userPassword model.UserChangePasswordBody
	if err := ctx.ShouldBindJSON(&userPassword); err != nil {
		middleware.HandleError(ctx, errors.NewValidationError("Неверный формат данных", err))
		return
	}
	if userPassword.Password != userPassword.PasswordConfirm {
		middleware.HandleError(ctx, errors.NewValidationError("Ошибка подтверждения пароля", err))
		return
	}

	success, err := h.Services.User.ChangePassword(id, userPassword)
	if err != nil {
		logger.GetLogger().Error("failed to change user password",
			zap.Error(err),
			zap.Int("user_id", id),
		)
		if strings.Contains(err.Error(), "пользователь не найден") {
			middleware.HandleError(ctx, errors.NewNotFoundError("пользователь", err))
			return
		}
		middleware.HandleError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, success)
}

// UserList
// @Summary Список пользователей
// @Tags Users
// @Produce		json
// @Success 200 {object} []model.User
// @Failure 401 {object} model.Error "Unauthorized"
// @Failure 500 {object} model.Error "Internal"
// @Router /api/v1/users [get]
// @Security BearerAuth.
func (h *Handler) ListUsers(ctx *gin.Context) {
	users, err := h.Services.User.GetAll()
	if err != nil {
		logger.GetLogger().Error("failed to get user list",
			zap.Error(err),
		)
		middleware.HandleError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, users)
}

// GetUserById
// @Summary Получение пользователя по id
// @Tags Users
// @Produce		json
// @Success 200 {object} model.User
// @Failure 400 {object} model.Error "Invalid request"
// @Failure 401 {object} model.Error "Unauthorized"
// @Failure 404 {object} model.Error "Not found"
// @Failure 500 {object} model.Error "Internal"
// @Param id path string true "id пользователя"
// @Router /api/v1/users/{id} [get]
// @Security BearerAuth.
func (h *Handler) GetUserByID(ctx *gin.Context) {
	idStr := ctx.Params.ByName("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		middleware.HandleError(ctx, errors.NewValidationError("Некорректный ID пользователя", err))
		return
	}
	user, err := h.Services.User.GetByID(id)
	if err != nil {
		logger.GetLogger().Error("failed to get user",
			zap.Error(err),
			zap.Int("user_id", id),
		)
		if strings.Contains(err.Error(), "пользователь не найден") {
			middleware.HandleError(ctx, errors.NewNotFoundError("пользователь", err))
			return
		}
		middleware.HandleError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, user)
}

// DeleteUser
// @Summary Удаление пользователя
// @Tags Users
// @Produce		json
// @Success 200 {object} model.Success "Объект успешно удален"
// @Failure 400 {object} model.Error "Invalid request"
// @Failure 401 {object} model.Error "Unauthorized"
// @Failure 404 {object} model.Error "Not found"
// @Failure 500 {object} model.Error "Internal"
// @Param id path string true "id пользователя"
// @Router /api/v1/users/{id} [delete]
// @Security BearerAuth.
func (h *Handler) DeleteUser(ctx *gin.Context) {
	role, exists := ctx.Get(userRoleKey)
	if !exists {
		middleware.HandleError(ctx, errors.NewValidationError("Некорректная роль пользователя", nil))
		return
	}
	isEmployee := role == model.RoleEmployee
	if !isEmployee {
		middleware.HandleError(ctx, errors.NewForbiddenError("Недостаточно прав для выполнения операции", nil))
		return
	}
	idStr := ctx.Params.ByName("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		middleware.HandleError(ctx, errors.NewValidationError("Некорректный ID пользователя", err))
		return
	}
	if err := h.Services.User.Delete(id); err != nil {
		logger.GetLogger().Error("failed to delete user",
			zap.Error(err),
			zap.Int("user_id", id),
		)
		if strings.Contains(err.Error(), "пользователь не найден") {
			middleware.HandleError(ctx, errors.NewNotFoundError("пользователь", err))
			return
		}
		middleware.HandleError(ctx, err)
		return
	}
	success := model.Success{
		Status:  "Success",
		Message: "Объект успешно удален",
	}
	ctx.JSON(http.StatusOK, success)
}
