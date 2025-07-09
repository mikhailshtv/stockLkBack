package model

import (
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type UserRole string

const (
	RoleClient   UserRole = "client"
	RoleEmployee UserRole = "employee"
)

type User struct {
	ID           int      `json:"id" db:"id"`
	Login        string   `json:"login" binding:"required" db:"login"`
	PasswordHash string   `json:"-" db:"password_hash"`
	FirstName    string   `json:"firstName" binding:"required" db:"first_name"`
	LastName     string   `json:"lastName" binding:"required" db:"last_name"`
	Email        string   `json:"email" binding:"required,email" db:"email"`
	Role         UserRole `json:"role" db:"role"`
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
	OldPassword     string `json:"oldPassword" binding:"required"`
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
	err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	return err == nil
}

func (user *User) HashPassword(password string) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	user.PasswordHash = string(bytes)
	return err
}

func (r UserRole) Valid() bool {
	switch r {
	case RoleClient, RoleEmployee:
		return true
	default:
		return false
	}
}
