package handler

import (
	"bytes"
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

	user.HashPassword(userProxy.Password)
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
			var userProxy *model.UserProxy
			userProxy, err = userProxy.UnmarshalJSONToUserProxy(reqBody)
			if err != nil {
				log.Fatalf("Ошибка десериализации: %v\n", err.Error())
			}
			v.HashPassword(userProxy.Password)
			if userProxy.Role == 0 {
				v.Role = 1
			} else {
				v.Role = userProxy.Role
			}
			ctx.Request.Body = io.NopCloser(bytes.NewBuffer(reqBody))
			if err := ctx.ShouldBindJSON(&v); err != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
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
