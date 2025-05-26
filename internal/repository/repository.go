package repository

import (
	"encoding/json"
	"errors"
	"golang/stockLkBack/internal/model"
	"io"
	"log"
	"os"
	"strings"
	"sync"
)

type Entity[T model.Order | model.Product | model.User] struct {
	Mu          sync.RWMutex
	Entities    []*T
	EntitiesLen int
}

func (entity *Entity[T]) AppendEntity(v T) {
	entity.Mu.Lock()
	defer entity.Mu.Unlock()
	entity.Entities = append(entity.Entities, &v)
}

func (entity *Entity[T]) SaveToFile(path string) {
	entity.Mu.Lock()
	defer entity.Mu.Unlock()
	outputPath := "./assets"
	if _, err := os.Stat(outputPath); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(outputPath, os.ModePerm)
		if err != nil {
			log.Fatalf("Ошибка создания каталога: %v\n", err.Error())
		}
	}
	if entity.EntitiesLen > 0 {
		switch any(entity.Entities[0]).(type) {
		case *model.User:
			usersJson := MarshalUserEntities(any(entity.Entities).([]*model.User))
			if err := os.WriteFile(path, usersJson, os.ModePerm); err != nil {
				log.Fatalf("Ошибка записи в файл: %v\n", err.Error())
			}
		default:
			json, err := json.Marshal(entity.Entities)
			if err != nil {
				log.Fatalf("Ошибка конвертирования в json: %v\n", err.Error())
			}
			if err := os.WriteFile(path, json, os.ModePerm); err != nil {
				log.Fatalf("Ошибка записи в файл: %v\n", err.Error())
			}
		}
	} else {
		json, err := json.Marshal(entity.Entities)
		if err != nil {
			log.Fatalf("Ошибка конвертирования в json: %v\n", err.Error())
		}
		if err := os.WriteFile(path, json, os.ModePerm); err != nil {
			log.Fatalf("Ошибка записи в файл: %v\n", err.Error())
		}
	}
}

func (entity *Entity[T]) RestoreFromFile(path string) {
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		return
	}
	entity.Mu.Lock()
	defer entity.Mu.Unlock()
	file, err := os.Open(path)
	if err != nil {
		log.Fatalf("Ошибка открытия файла: %v\n", err.Error())
	}
	defer file.Close()
	data, err := io.ReadAll(file)
	if err != nil {
		log.Fatalf("Ошибка чтения из файла: %v\n", err.Error())
	}
	if len(data) == 0 {
		return
	}

	if strings.Contains(path, "users") {
		if err := UnmarshalingUserEntitiesJson(data); err != nil {
			log.Fatalf("Ошибка десериализации: %v\n", err.Error())
		}
	} else {
		jsonError := json.Unmarshal(data, &entity.Entities)
		entity.EntitiesLen = len(entity.Entities)
		if jsonError != nil {
			log.Fatalf("Ошибка десериализации: %v\n", jsonError.Error())
		}
	}
}

func (entity *Entity[T]) SavedEntities() []*T {
	entity.Mu.RLock()
	defer entity.Mu.RUnlock()
	return entity.Entities[entity.EntitiesLen:]
}

var OrdersStruct = Entity[model.Order]{}
var ProductsStruct = Entity[model.Product]{}
var UsersStruct = Entity[model.User]{}

func CheckAndSaveEntity(entity any) {
	switch v := entity.(type) {
	case model.Order:
		OrdersStruct.AppendEntity(v)
		OrdersStruct.SaveToFile("./assets/orders.json")
	case model.Product:
		ProductsStruct.AppendEntity(v)
		ProductsStruct.SaveToFile("./assets/products.json")
	case model.User:
		UsersStruct.AppendEntity(v)
		UsersStruct.SaveToFile("./assets/users.json")
	}
}

func MarshalUserEntities(entities []*model.User) []byte {
	proxyUsers := make([]struct {
		Id           int            `json:"id"`
		Login        string         `json:"login"`
		PasswordHash string         `json:"password"`
		FirstName    string         `json:"firstName"`
		LastName     string         `json:"lastName"`
		Email        string         `json:"email"`
		Role         model.UserRole `json:"role"`
	}, len(entities), cap(entities))

	for i, item := range entities {
		proxyUsers[i].Id = item.Id
		proxyUsers[i].Login = item.Login
		proxyUsers[i].PasswordHash = item.PasswordHash()
		proxyUsers[i].FirstName = item.FirstName
		proxyUsers[i].LastName = item.LastName
		proxyUsers[i].Email = item.Email
		proxyUsers[i].Role = item.Role
	}
	val, err := json.Marshal(proxyUsers)
	if err != nil {
		log.Fatalf("Ошибка конвертирования в json: %v\n", err.Error())
	}
	return val
}

func UnmarshalingUserEntitiesJson(data []byte) error {
	var temp []struct {
		Id           int            `json:"id"`
		Login        string         `json:"login"`
		PasswordHash string         `json:"password"`
		FirstName    string         `json:"firstName"`
		LastName     string         `json:"lastName"`
		Email        string         `json:"email"`
		Role         model.UserRole `json:"role"`
	}

	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}

	for _, v := range temp {
		currentUser := &model.User{
			Id:        v.Id,
			Login:     v.Login,
			FirstName: v.FirstName,
			LastName:  v.LastName,
			Email:     v.Email,
			Role:      v.Role,
		}
		currentUser.SetPasswordHash(v.PasswordHash)
		UsersStruct.Entities = append(UsersStruct.Entities, currentUser)
	}
	return nil
}
