package service

import (
	"context"
	"encoding/json"
	"errors"
	"golang/stockLkBack/internal/model"
	"golang/stockLkBack/internal/repository"
	"log"
	"sync"
	"time"

	"github.com/ddosify/go-faker/faker"
)

func NewEntity() any {
	caseNumber := faker.NewFaker().RandomIntBetween(0, 2)
	switch caseNumber {
	case 0:
		return NewOrder()
	case 1:
		return NewProduct()
	case 2:
		return NewUser()
	default:
		return errors.New("ошибка при создании сущности")
	}
}

func NewOrder() model.Order {
	order := model.Order{
		Id:          int(faker.NewFaker().RandomUUID().ID()),
		Number:      faker.NewFaker().RandomIntBetween(1, 999999),
		TotalCost:   faker.NewFaker().RandomIntBetween(1, 999999),
		CreatedDate: time.Now().Add(time.Hour * time.Duration(faker.NewFaker().RandomIntBetween(-8760, -1))),
		Status:      model.OrderStatus(faker.NewFaker().RandomIntBetween(1, 2)),
	}
	// order.SetLastModifiedDate(time.Now().Add(time.Hour * time.Duration(faker.NewFaker().RandomIntBetween(-8760, -1))))
	return order
}

func NewProduct() model.Product {
	product := model.Product{
		Id:        int(faker.NewFaker().RandomUUID().ID()),
		Code:      faker.NewFaker().RandomIntBetween(1, 999999),
		Quantity:  faker.NewFaker().RandomIntBetween(1, 9999),
		Name:      faker.NewFaker().RandomProduct(),
		SalePrice: faker.NewFaker().RandomIntBetween(1, 999999),
	}
	// product.SetPurchasePrice(faker.NewFaker().RandomIntBetween(1, 999999))
	return product
}

func NewUser() model.User {
	user := model.User{
		Id:        int(faker.NewFaker().RandomUUID().ID()),
		Login:     faker.NewFaker().RandomUsername(),
		FirstName: faker.NewFaker().RandomPersonFirstName(),
		LastName:  faker.NewFaker().RandomPersonLastName(),
		Email:     faker.NewFaker().RandomEmail(),
		Role:      model.UserRole(faker.NewFaker().RandomIntBetween(0, 1)),
	}
	err := user.HashPassword(faker.NewFaker().RandomPassword())
	if err != nil {
		log.Fatal(err.Error())
	}
	return user
}

func LogAddedEntities(ctx context.Context) {
	go func() {
		for {
			t := time.NewTicker(time.Millisecond * 200)
			select {
			case <-ctx.Done():
				log.Printf("Received %v signal, shutting down log goroutine \n", ctx.Err().Error())
				return
			case <-t.C:
				ordersJSON, err := json.Marshal(repository.OrdersStruct.SavedEntities())
				if err != nil {
					log.Fatal(err.Error())
				} else if len(repository.OrdersStruct.SavedEntities()) > 0 {
					log.Printf("Orders: %v\n", string(ordersJSON))
				}
				productsJSON, err := json.Marshal(repository.ProductsStruct.SavedEntities())
				if err != nil {
					log.Fatal(err.Error())
				} else if len(repository.ProductsStruct.SavedEntities()) > 0 {
					log.Printf("Products: %v\n", string(productsJSON))
				}
				usersJSON, err := json.Marshal(repository.UsersStruct.SavedEntities())
				if err != nil {
					log.Fatal(err.Error())
				} else if len(repository.UsersStruct.SavedEntities()) > 0 {
					log.Printf("Users: %v\n", string(usersJSON))
				}

				repository.OrdersStruct.EntitiesLen = len(repository.OrdersStruct.Entities)
				repository.ProductsStruct.EntitiesLen = len(repository.ProductsStruct.Entities)
				repository.UsersStruct.EntitiesLen = len(repository.UsersStruct.Entities)
			}
		}
	}()
}

func Interval(ctx context.Context, wg *sync.WaitGroup) {
	channel := make(chan any, 1)
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			t := time.NewTicker(time.Second)
			select {
			case <-ctx.Done():
				log.Printf("Received %v signal, shutting down creating entity goroutine \n", ctx.Err().Error())
				return
			case <-t.C:
				channel <- NewEntity()
			}
		}
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-ctx.Done():
				log.Printf("Received %v signal, shutting down check and save entity goroutine \n", ctx.Err().Error())
				return
			case entity := <-channel:
				repository.CheckAndSaveEntity(entity)
			}
		}
	}()
}
