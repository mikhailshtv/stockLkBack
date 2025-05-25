package model

import (
	"encoding/json"

	"golang.org/x/crypto/bcrypt"
)

type UserRole int

const (
	Client UserRole = iota
	Employee
)

type User struct {
	Id           int      `json:"id"`
	Login        string   `json:"login" binding:"required"`
	passwordHash string   `json:"-"`
	FirstName    string   `json:"firstName" binding:"required"`
	LastName     string   `json:"lastName" binding:"required"`
	Email        string   `json:"email" binding:"required"`
	Role         UserRole `json:"role"`
}

func (user *User) CheckUserPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(user.passwordHash), []byte(password))
	return err == nil
}

func (user *User) HashPassword(password string) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	user.passwordHash = string(bytes)
	return err
}
func (user *User) PasswordHash() string {
	return user.passwordHash
}

func (user *User) SetPasswordHash(passwordHash string) {
	user.passwordHash = passwordHash
}

func (user *User) UnmarshalUserJSON(data []byte) error {
	var temp struct {
		Id           int      `json:"id"`
		Login        string   `json:"login"`
		PasswordHash string   `json:"password"`
		FirstName    string   `json:"firstName"`
		LastName     string   `json:"lastName"`
		Email        string   `json:"email"`
		Role         UserRole `json:"role"`
	}

	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}
	user.Id = temp.Id
	user.Login = temp.Login
	user.passwordHash = temp.PasswordHash
	user.FirstName = temp.FirstName
	user.LastName = temp.LastName
	user.Email = temp.Email
	user.Role = temp.Role
	return nil
}
