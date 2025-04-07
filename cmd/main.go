package main

import "golang/stockLkBack/internal/model"

func main() {
	user := model.User{}
	user.SetPasswordHash("Qwerty123")
}
