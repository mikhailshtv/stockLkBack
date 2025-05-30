package handler

import (
	"bytes"
	"golang/stockLkBack/internal/model"
	"golang/stockLkBack/internal/repository"
	"golang/stockLkBack/internal/utils/jwtgen"
	"io"
	"log"
	"net/http"
	"slices"
	"strconv"

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
// @Router /api/v1/users [post]
// @Security BearerAuth
func CreateUser(ctx *gin.Context) {
	reqBody, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		log.Fatal(err)
	}
	var userProxy *model.UserProxy
	userProxy, err = userProxy.UnmarshalJSONToUserProxy(reqBody)
	if err != nil {
		log.Fatalf("Ошибка десериализации: %v\n", err.Error())
	}
	var user model.User
	if repository.UsersStruct.EntitiesLen == 0 {
		user.Id = 1
	} else {
		length := repository.UsersStruct.EntitiesLen
		user.Id = repository.UsersStruct.Entities[length-1].Id + 1
	}

	if userProxy.Password == userProxy.PasswordConfirm {
		user.HashPassword(userProxy.Password)
	} else {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Ошибка подтверждения пароля"})
		return
	}

	if userProxy.Role == 0 {
		user.Role = 1
	}
	ctx.Request.Body = io.NopCloser(bytes.NewBuffer(reqBody))
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	repository.CheckAndSaveEntity(user)
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
// @Router /api/v1/login [post]
func Login(ctx *gin.Context) {
	var loginRequest model.LoginRequest
	if err := ctx.ShouldBindJSON(&loginRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		ctx.Abort()
		return
	}

	var index int
	for i, v := range repository.UsersStruct.Entities {
		if v.Login == loginRequest.Login && v.CheckUserPassword(loginRequest.Password) {
			index = i
			// Генерация токена
			token, err := jwtgen.GenerateToken(loginRequest.Login, v.Role)
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка генерации токена"})
				return
			}

			ctx.JSON(http.StatusOK, gin.H{
				"message": "Login successful",
				"token":   token,
			})
		}
	}
	if index == 0 {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Логин или пароль пользователя недействителен"})
		return
	}

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
// @Security BearerAuth
func EditUser(ctx *gin.Context) {
	idStr := ctx.Params.ByName("id")
	var index int
	for i, v := range repository.UsersStruct.Entities {
		id, err := strconv.Atoi(idStr)
		if err != nil {
			log.Fatal(err)
		}
		if v.Id == id {
			index = i
			var userEdit *model.UserEditBody
			if err := ctx.ShouldBindJSON(&userEdit); err != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			if userEdit.Email == "" && userEdit.FirstName == "" && userEdit.LastName == "" {
				ctx.JSON(http.StatusBadRequest, gin.H{"error": "Недопустимое тело запроса"})
				ctx.Abort()
				return
			}
			if userEdit.Email != "" {
				v.Email = userEdit.Email
			}
			if userEdit.FirstName != "" {
				v.FirstName = userEdit.FirstName
			}
			if userEdit.LastName != "" {
				v.LastName = userEdit.LastName
			}
			repository.UsersStruct.SaveToFile("./assets/users.json")
			ctx.JSON(http.StatusOK, v)
		}
	}
	if index == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Объект не найден"})
		ctx.Abort()
		return
	}
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
// @Security BearerAuth
func ChangeUserRole(ctx *gin.Context) {
	idStr := ctx.Params.ByName("id")
	for _, v := range repository.UsersStruct.Entities {
		id, err := strconv.Atoi(idStr)
		if err != nil {
			log.Fatal(err)
		}
		if v.Id == id {
			var userRole *model.UserRoleBody
			if err := ctx.ShouldBindJSON(&userRole); err != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			if userRole.Role == 0 || userRole.Role > 2 {
				ctx.JSON(http.StatusBadRequest, gin.H{"error": "Недопустимое тело запроса"})
				return
			} else {
				v.Role = userRole.Role
			}
			repository.UsersStruct.SaveToFile("./assets/users.json")
			ctx.JSON(http.StatusOK, v)
		}
	}
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
// @Security BearerAuth
func ChangeUserPassword(ctx *gin.Context) {
	idStr := ctx.Params.ByName("id")
	for _, v := range repository.UsersStruct.Entities {
		id, err := strconv.Atoi(idStr)
		if err != nil {
			log.Fatal(err)
		}
		if v.Id == id {
			var userPassword *model.UserChangePasswordBody
			if err := ctx.ShouldBindJSON(&userPassword); err != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			if userPassword.Password != userPassword.PasswordConfirm {
				ctx.JSON(http.StatusBadRequest, gin.H{"error": "Ошибка подтверждения пароля"})
				return
			} else {
				v.HashPassword(userPassword.Password)
			}
			repository.UsersStruct.SaveToFile("./assets/users.json")
			success := model.Success{
				Status:  "Success",
				Message: "Пароль успешно изменен",
			}
			ctx.JSON(http.StatusOK, success)
		}
	}
}

// UserList
// @Summary Список пользователей
// @Tags Users
// @Produce		json
// @Success 200 {object} []model.User
// @Failure 400 {string} string "Invalid request"
// @Router /api/v1/users [get]
// @Security BearerAuth
func ListUsers(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, repository.UsersStruct.Entities)
}

// GetUserById
// @Summary Получение пользователя по id
// @Tags Users
// @Produce		json
// @Success 200 {object} model.User
// @Failure 400 {string} string "Invalid request"
// @Param id path string true "id пользователя"
// @Router /api/v1/users/{id} [get]
// @Security BearerAuth
func GetUserById(ctx *gin.Context) {
	idStr := ctx.Params.ByName("id")
	for _, v := range repository.UsersStruct.Entities {
		id, err := strconv.Atoi(idStr)
		if err != nil {
			log.Fatal(err)
		}
		if v.Id == id {
			ctx.JSON(http.StatusOK, v)
		}
	}
}

// DeleteUser
// @Summary Удаление пользователя
// @Tags Users
// @Produce		json
// @Success 200 {object} model.Success "Объект успешно удален"
// @Failure 400 {string} string "Invalid request"
// @Param id path string true "id пользователя"
// @Router /api/v1/users/{id} [delete]
// @Security BearerAuth
func DeleteUser(ctx *gin.Context) {
	idStr := ctx.Params.ByName("id")
	for i, v := range repository.UsersStruct.Entities {
		id, err := strconv.Atoi(idStr)
		if err != nil {
			log.Fatal(err)
		}
		if v.Id == id {
			repository.UsersStruct.Entities = slices.Delete(repository.UsersStruct.Entities, i, i+1)
			repository.UsersStruct.EntitiesLen = len(repository.UsersStruct.Entities)
			repository.UsersStruct.SaveToFile("./assets/users.json")
			success := model.Success{
				Status:  "Success",
				Message: "Объект успешно удален",
			}
			ctx.JSON(http.StatusOK, success)
		}
	}
}
