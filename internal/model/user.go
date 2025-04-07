package model

import (
	"log"

	"golang.org/x/crypto/bcrypt"
)

type UserRole int

const (
	Client UserRole = iota
	Employee
)

type User struct {
	Id           int
	Login        string
	passwordHash string
	FirstName    string
	LastName     string
	Email        string
	Role         UserRole
}

func (user *User) CheckUserPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(user.passwordHash), []byte(password))
	return err == nil
}

func (user *User) SetPasswordHash(password string) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 32)
	if err != nil {
		log.Fatal(err.Error())
	} else {
		user.passwordHash = string(bytes)
	}
}
