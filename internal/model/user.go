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

type UserProxy struct {
	Id        int      `json:"id"`
	Login     string   `json:"login"`
	Password  string   `json:"password"`
	FirstName string   `json:"firstName"`
	LastName  string   `json:"lastName"`
	Email     string   `json:"email"`
	Role      UserRole `json:"role,omitempty"`
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
