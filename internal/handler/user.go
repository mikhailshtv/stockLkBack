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
			var user model.User
			reqBody, err := io.ReadAll(ctx.Request.Body)
			if err != nil {
				log.Fatal(err)
			}
			reqBodyMap := make(map[string]any)
			if err := json.Unmarshal(reqBody, &reqBodyMap); err != nil {
				log.Fatal(err)
			}
			ctx.Request.Body = io.NopCloser(bytes.NewBuffer(reqBody))
			if err := ctx.ShouldBindJSON(&user); err != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			v.Login = user.Login
			v.HashPassword(reqBodyMap["password"].(string))
			v.FirstName = user.FirstName
			v.LastName = user.LastName
			v.Email = user.Email
			if reqBodyMap["role"] == nil {
				v.Role = 1
			} else {
				v.Role = user.Role
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
