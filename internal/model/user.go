package model

import (
	"encoding/json"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type UserRole int

const (
	Client UserRole = iota + 1
	Employee
)

type User struct {
	ID           int      `json:"id"`
	Login        string   `json:"login" binding:"required"`
	passwordHash string   `json:"-"`
	FirstName    string   `json:"firstName" binding:"required"`
	LastName     string   `json:"lastName" binding:"required"`
	Email        string   `json:"email" binding:"required,email"`
	Role         UserRole `json:"role"`
}

type UserProxy struct {
	ID              int      `json:"id"`
	Login           string   `json:"login"`
	Password        string   `json:"password"`
	PasswordConfirm string   `json:"passwordConfirm"`
	FirstName       string   `json:"firstName"`
	LastName        string   `json:"lastName"`
	Email           string   `json:"email"`
	Role            UserRole `json:"role,omitempty"`
}

type UserCreateBody struct {
	Login           string `json:"login"`
	Password        string `json:"password"`
	PasswordConfirm string `json:"passwordConfirm"`
	FirstName       string `json:"firstName"`
	LastName        string `json:"lastName"`
	Email           string `json:"email"`
}

type UserEditBody struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
}

type UserRoleBody struct {
	Role UserRole `json:"role" binding:"required"`
}

type UserChangePasswordBody struct {
	Password        string `json:"password" binding:"required"`
	PasswordConfirm string `json:"passwordConfirm" binding:"required"`
}

type LoginRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

// Claims Структура для JWT токена.
type Claims struct {
	Login                string   `json:"login"`
	Role                 UserRole `json:"role"`
	jwt.RegisteredClaims          // Данное поле нужно для правильной генерации JWT.
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

func (userProxy *UserProxy) UnmarshalJSONToUserProxy(data []byte) (*UserProxy, error) {
	if err := json.Unmarshal(data, &userProxy); err != nil {
		return nil, err
	}
	return userProxy, nil
}
