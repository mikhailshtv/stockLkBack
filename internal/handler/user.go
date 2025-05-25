package handler

import (
	"bytes"
	"encoding/json"
	"golang/stockLkBack/internal/model"
	"golang/stockLkBack/internal/repository"
	"io"
	"log"
	"net/http"

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
	log.Println(string(reqBody))
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
