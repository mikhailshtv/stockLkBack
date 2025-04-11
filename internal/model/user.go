package model

import (
	"golang.org/x/crypto/bcrypt"
)

type UserRole int

const (
	Client UserRole = iota
	Employee
)

type User struct {
	Id           uint32
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

func (user *User) SetPasswordHash(password string) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	user.passwordHash = string(bytes)
	return err
}
