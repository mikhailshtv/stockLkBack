package repository

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"os"
	"slices"

	"golang/stockLkBack/internal/model"
	"golang/stockLkBack/internal/utils/jwtgen"

	"go.mongodb.org/mongo-driver/mongo"
)

type UsersRepository struct {
	Users    []model.User
	UsersLen int
	db       *mongo.Database
}

func NewUsersRepository(db *mongo.Database) *UsersRepository {
	return &UsersRepository{db: db}
}

func (ur *UsersRepository) Create(userRequest model.UserCreateBody) (*model.User, error) {
	var user model.User
	if ur.UsersLen > 0 {
		lastUser := ur.Users[ur.UsersLen-1]
		user.ID = lastUser.ID + 1
	} else {
		user.ID = 1
	}
	user.FirstName = userRequest.FirstName
	user.LastName = userRequest.LastName
	user.Login = userRequest.Login
	user.Email = userRequest.Email
	if userRequest.Password == userRequest.PasswordConfirm {
		err := user.HashPassword(userRequest.Password)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, errors.New("ошибка подтверждения пароля")
	}

	ur.Users = append(ur.Users, user)
	ur.UsersLen = len(ur.Users)

	if err := saveUsersToFile(ur.Users); err != nil {
		return nil, err
	}
	return &user, nil
}

func (ur *UsersRepository) GetAll() ([]model.User, error) {
	return ur.Users, nil
}

func (ur *UsersRepository) GetByID(id int) (*model.User, error) {
	idx := slices.IndexFunc(ur.Users, func(user model.User) bool { return user.ID == id })
	if idx == -1 {
		return nil, errors.New(NotFoundErrorMessage)
	}
	return &ur.Users[idx], nil
}

func (ur *UsersRepository) Delete(id int) error {
	ur.Users = slices.DeleteFunc(ur.Users, func(user model.User) bool { return user.ID == id })
	ur.UsersLen = len(ur.Users)
	if err := saveUsersToFile(ur.Users); err != nil {
		return err
	}
	return nil
}

func (ur *UsersRepository) Update(id int, userReq model.UserEditBody) (*model.User, error) {
	idx := slices.IndexFunc(ur.Users, func(user model.User) bool { return user.ID == id })
	if idx == -1 {
		return nil, errors.New(NotFoundErrorMessage)
	}
	if userReq.FirstName != "" {
		ur.Users[idx].FirstName = userReq.FirstName
	}
	if userReq.LastName != "" {
		ur.Users[idx].LastName = userReq.LastName
	}
	if userReq.Email != "" {
		ur.Users[idx].Email = userReq.Email
	}

	if err := saveUsersToFile(ur.Users); err != nil {
		return nil, err
	}
	return &ur.Users[idx], nil
}

func (ur *UsersRepository) Login(userReq model.LoginRequest) (*model.TokenSuccess, error) {
	idx := slices.IndexFunc(ur.Users, func(user model.User) bool { return user.Login == userReq.Login })
	if idx == -1 {
		return nil, errors.New("логин или пароль пользователя недействителен")
	}
	foundUser := ur.Users[idx]
	if foundUser.CheckUserPassword(userReq.Password) {
		token, err := jwtgen.GenerateToken(userReq.Login, foundUser.Role)
		if err != nil {
			return nil, errors.New("ошибка генерации токена")
		}
		return &model.TokenSuccess{
			Message: "Login successful",
			Token:   token,
		}, nil
	}
	return nil, errors.New("логин или пароль пользователя недействителен")
}

func (ur *UsersRepository) ChangeUserRole(id int, userRoleReq model.UserRoleBody) (*model.User, error) {
	idx := slices.IndexFunc(ur.Users, func(user model.User) bool { return user.ID == id })
	if idx == -1 {
		return nil, errors.New("пользователь не найден")
	}
	foundUser := ur.Users[idx]
	foundUser.Role = userRoleReq.Role
	return &foundUser, nil
}

func (ur *UsersRepository) ChangePassword(
	id int,
	changePassworReq model.UserChangePasswordBody,
) (*model.Success, error) {
	idx := slices.IndexFunc(ur.Users, func(user model.User) bool { return user.ID == id })
	if idx == -1 {
		return nil, errors.New("пользователь не найден")
	}
	foundUser := ur.Users[idx]
	err := foundUser.HashPassword(changePassworReq.Password)
	if err != nil {
		return nil, err
	}
	return &model.Success{
		Status:  "Success",
		Message: "Пароль успешно изменен",
	}, nil
}

func (ur *UsersRepository) RestoreUsersFromFile(path string) {
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		return
	}
	file, err := os.Open(path)
	if err != nil {
		log.Printf("Ошибка открытия файла: %v\n", err.Error())
		return
	}
	defer file.Close()
	data, err := io.ReadAll(file)
	if err != nil {
		log.Printf("Ошибка чтения из файла: %v\n", err.Error())
		return
	}
	if len(data) == 0 {
		return
	}

	users, err := UnmarshalingUserEntitiesJSON(data)
	if err != nil {
		log.Printf("Ошибка десериализации: %v\n", err.Error())
		return
	}
	ur.Users = users
}

func saveUsersToFile(users []model.User) error {
	outputPath := "./assets"
	if _, err := os.Stat(outputPath); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(outputPath, os.ModePerm)
		if err != nil {
			log.Printf("Ошибка создания каталога: %v\n", err.Error())
			return err
		}
	}
	path := "./assets/users.json"
	json, err := json.Marshal(users)
	if err != nil {
		log.Printf("Ошибка конвертирования в json: %v\n", err.Error())
		return err
	}
	if err := os.WriteFile(path, json, 0o600); err != nil {
		log.Printf("Ошибка записи в файл: %v\n", err.Error())
		return err
	}
	return nil
}

func UnmarshalingUserEntitiesJSON(data []byte) ([]model.User, error) {
	var temp []struct {
		ID           int            `json:"id"`
		Login        string         `json:"login"`
		PasswordHash string         `json:"password"`
		FirstName    string         `json:"firstName"`
		LastName     string         `json:"lastName"`
		Email        string         `json:"email"`
		Role         model.UserRole `json:"role"`
	}

	if err := json.Unmarshal(data, &temp); err != nil {
		return nil, err
	}

	users := make([]model.User, 0, 10)

	for _, v := range temp {
		currentUser := model.User{
			ID:        v.ID,
			Login:     v.Login,
			FirstName: v.FirstName,
			LastName:  v.LastName,
			Email:     v.Email,
			Role:      v.Role,
		}
		currentUser.SetPasswordHash(v.PasswordHash)
		users = append(users, currentUser)
	}
	return users, nil
}
