package service

import (
	"golang/stockLkBack/internal/repository"
)

func RestoreData() {
	repository.OrdersStruct.RestoreFromFile("./assets/orders.json")
	// ordersJSON, err := json.Marshal(repository.OrdersStruct.Entities)
	// if err != nil {
	// 	log.Fatalf("Ошибка сериализации %v", err.Error())
	// }
	// log.Printf("восстановленные данные orders: %v\n", string(ordersJSON))

	repository.ProductsStruct.RestoreFromFile("./assets/products.json")
	// productsJSON, err := json.Marshal(repository.ProductsStruct.Entities)
	// if err != nil {
	// 	log.Fatalf("Ошибка сериализации %v", err.Error())
	// }
	// log.Printf("восстановленные данные products: %v\n", string(productsJSON))

	repository.UsersStruct.RestoreFromFile("./assets/users.json")
	// usersJSON, err := json.Marshal(repository.UsersStruct.Entities)
	// if err != nil {
	// 	log.Fatalf("Ошибка сериализации %v", err.Error())
	// }
	// log.Printf("восстановленные данные users: %v\n", string(usersJSON))
}
