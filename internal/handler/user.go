package handler

import (
	"bytes"
	"encoding/json"
	"golang/stockLkBack/internal/model"
	"golang/stockLkBack/internal/repository"
	"io"
	"log"
	"net/http"
	"slices"
	"strconv"

	"github.com/gin-gonic/gin"
)

func CreateUser(ctx *gin.Context) {
	var user model.User
	user.Id = repository.UsersStruct.EntitiesLen + 1
	reqBody, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		log.Fatal(err)
	}
	reqBodyMap := make(map[string]any)
	if err := json.Unmarshal(reqBody, &reqBodyMap); err != nil {
		log.Fatal(err)
	}
	user.HashPassword(reqBodyMap["password"].(string))
	if reqBodyMap["role"] == nil {
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

func EditUser(ctx *gin.Context) {
	idStr := ctx.Params.ByName("id")
	for _, v := range repository.UsersStruct.Entities {
		id, err := strconv.Atoi(idStr)
		if err != nil {
			log.Fatal(err)
		}
		if v.Id == id {
			reqBody, err := io.ReadAll(ctx.Request.Body)
			if err != nil {
				log.Fatal(err)
			}
			var userProxy struct {
				Id        int            `json:"id"`
				Login     string         `json:"login"`
				Password  string         `json:"password"`
				FirstName string         `json:"firstName"`
				LastName  string         `json:"lastName"`
				Email     string         `json:"email"`
				Role      model.UserRole `json:"role,omitempty"`
			}
			if err := json.Unmarshal(reqBody, &userProxy); err != nil {
				log.Fatal(err)
			}
			v.Login = userProxy.Login
			v.HashPassword(userProxy.Password)
			v.FirstName = userProxy.FirstName
			v.LastName = userProxy.LastName
			v.Email = userProxy.Email
			v.Role = userProxy.Role
			if userProxy.Role == 0 {
				v.Role = 1
			} else {
				v.Role = userProxy.Role
			}
			repository.UsersStruct.SaveToFile("./assets/users.json")
			ctx.JSON(http.StatusOK, v)
		}
	}
}

func ListUsers(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, repository.UsersStruct.Entities)
}

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
			ctx.JSON(http.StatusOK, gin.H{"status": "Success", "message": "Объект успешно удален"})
		}
	}
}
