package handler

import (
	"log"
	"net/http"
	"strconv"

	"golang/stockLkBack/internal/model"

	"github.com/gin-gonic/gin"
)

// CreateUser
// @Summary Создание/регистрация пользователя
// @Tags Users
// @Accept			json
// @Produce		json
// @Param user body model.UserCreateBody true "Объект пользователя"
// @Success 200 {object} model.User
// @Failure 400 {object} model.Error "Invalid request"
// @Failure 500 {object} model.Error "Internal"
// @Router /api/v1/users [post].
func (h *Handler) CreateUser(ctx *gin.Context) {
	var userReq model.UserCreateBody
	if err := ctx.ShouldBindJSON(&userReq); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user, err := h.Services.User.Create(userReq)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, user)
}

// Login
// @Summary Аутенификация пользователя
// @Tags Login
// @Accept			json
// @Produce		json
// @Param user body model.LoginRequest true "Данные для аутентификации пользователя"
// @Success 200 {object} model.TokenSuccess
// @Failure 400 {object} model.Error "Invalid request"
// @Failure 401 {object} model.Error "Anauthorized"
// @Failure 500 {object} model.Error "Internal server error"
// @Router /api/v1/login [post].
func (h *Handler) Login(ctx *gin.Context) {
	var loginRequest model.LoginRequest
	if err := ctx.ShouldBindJSON(&loginRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		ctx.Abort()
		return
	}
	TokenSuccess, err := h.Services.User.Login(loginRequest)
	if err != nil {
		if err.Error() == "ошибка генерации токена" {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		} else if err.Error() == "логин или пароль пользователя недействителен" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			ctx.Abort()
			return
		}
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
// @Param id path string true "id пользователя"
// @Router /api/v1/users/{id} [put]
// @Security BearerAuth.
func (h *Handler) EditUser(ctx *gin.Context) {
	idStr := ctx.Params.ByName("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Fatal(err)
	}
	var userEdit model.UserEditBody
	if err := ctx.ShouldBindJSON(&userEdit); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if userEdit.Email == "" && userEdit.FirstName == "" && userEdit.LastName == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Недопустимое тело запроса"})
		ctx.Abort()
		return
	}
	user, err := h.Services.User.Update(id, userEdit)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		ctx.Abort()
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
// @Param id path string true "id пользователя"
// @Router /api/v1/users/{id}/role [patch]
// @Security BearerAuth.
func (h *Handler) ChangeUserRole(ctx *gin.Context) {
	isEmployee := ctx.GetString(userRoleKey) == "employee"
	if !isEmployee {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "недостаточно прав для выполнения операции"})
		return
	}
	idStr := ctx.Params.ByName("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Fatal(err)
	}
	var userRole *model.UserRoleBody
	if err := ctx.ShouldBindJSON(&userRole); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if !userRole.Role.Valid() {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Недопустимое тело запроса"})
		return
	}

	user, err := h.Services.User.ChangeUserRole(id, *userRole)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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
// @Param id path string true "id пользователя"
// @Router /api/v1/users/{id}/password [patch]
// @Security BearerAuth.
func (h *Handler) ChangeUserPassword(ctx *gin.Context) {
	idStr := ctx.Params.ByName("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Fatal(err)
	}
	var userPassword model.UserChangePasswordBody
	if err := ctx.ShouldBindJSON(&userPassword); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if userPassword.Password != userPassword.PasswordConfirm {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Ошибка подтверждения пароля"})
		return
	}

	success, err := h.Services.User.ChangePassword(id, userPassword)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, success)
}

// UserList
// @Summary Список пользователей
// @Tags Users
// @Produce		json
// @Success 200 {object} []model.User
// @Failure 400 {string} string "Invalid request"
// @Router /api/v1/users [get]
// @Security BearerAuth.
func (h *Handler) ListUsers(ctx *gin.Context) {
	users, err := h.Services.User.GetAll()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		ctx.Abort()
		return
	}
	ctx.JSON(http.StatusOK, users)
}

// GetUserById
// @Summary Получение пользователя по id
// @Tags Users
// @Produce		json
// @Success 200 {object} model.User
// @Failure 400 {string} string "Invalid request"
// @Param id path string true "id пользователя"
// @Router /api/v1/users/{id} [get]
// @Security BearerAuth.
func (h *Handler) GetUserByID(ctx *gin.Context) {
	idStr := ctx.Params.ByName("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Fatal(err)
	}
	user, err := h.Services.User.GetByID(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		ctx.Abort()
		return
	}
	ctx.JSON(http.StatusOK, user)
}

// DeleteUser
// @Summary Удаление пользователя
// @Tags Users
// @Produce		json
// @Success 200 {object} model.Success "Объект успешно удален"
// @Failure 400 {string} string "Invalid request"
// @Param id path string true "id пользователя"
// @Router /api/v1/users/{id} [delete]
// @Security BearerAuth.
func (h *Handler) DeleteUser(ctx *gin.Context) {
	isEmployee := ctx.GetString(userRoleKey) == "employee"
	if !isEmployee {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "недостаточно прав для выполнения операции"})
		return
	}
	idStr := ctx.Params.ByName("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Fatal(err)
	}
	if err := h.Services.User.Delete(id); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		ctx.Abort()
		return
	}
	success := model.Success{
		Status:  "Success",
		Message: "Объект успешно удален",
	}
	ctx.JSON(http.StatusOK, success)
}
