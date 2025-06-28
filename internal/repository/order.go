package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"golang/stockLkBack/internal/model"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type OrdersRepository struct {
	db             *mongo.Database
	redis          *redis.Client
	collectionName string
}

func NewOrdersRepository(db *mongo.Database, redis *redis.Client, collectionName string) *OrdersRepository {
	return &OrdersRepository{db: db, redis: redis, collectionName: collectionName}
}

func (or *OrdersRepository) Create(orderRequest model.OrderRequestBody) (*model.Order, error) {
	var order model.Order
	ordersCollection := or.db.Collection(or.collectionName)
	_, err := ordersCollection.Find(context.TODO(), bson.D{})
	if err != nil {
		fmt.Println("Коллекция пуста")
		order.Id = 1
		order.Number = 1
	} else {
		orderNextId, err := getNextSequence(or.db, "orderid")
		if err != nil {
			return nil, err
		}
		orderNextNumber, err := getNextSequence(or.db, "orderNumber")
		if err != nil {
			return nil, err
		}
		order.Id = orderNextId
		order.Number = orderNextNumber
	}

	order.CreatedDate = time.Now().UTC()
	order.LastModifiedDate = time.Now().UTC()
	order.Status = model.Active
	totalCost := 0
	for _, product := range orderRequest.Products {
		totalCost += product.SalePrice
	}
	order.TotalCost = totalCost
	order.Products = orderRequest.Products
	_, err = ordersCollection.InsertOne(context.TODO(), order)
	if err != nil {
		log.Fatal(err)
	}

	return &order, nil
}

func (or *OrdersRepository) GetAll() ([]model.Order, error) {
	cursor, err := or.db.Collection(or.collectionName).Find(context.TODO(), bson.D{})
	if err != nil {
		return nil, err
	}
	var orders []model.Order
	if err = cursor.All(context.TODO(), &orders); err != nil {
		return nil, err
	}
	return orders, nil
}

func (or *OrdersRepository) GetById(id int) (*model.Order, error) {
	ordersCollection := or.db.Collection(or.collectionName)
	var result bson.M
	filter := bson.M{"_id": id}
	err := ordersCollection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		return nil, errors.New("элемент не найден")
	}
	var order model.Order
	resultBytes, err := bson.Marshal(result)
	if err != nil {
		return nil, err
	}
	bson.Unmarshal(resultBytes, &order)

	return &order, nil
}

func (or *OrdersRepository) Delete(id int) (*model.Order, error) {
	ordersCollection := or.db.Collection(or.collectionName)
	filter := bson.M{"_id": id}
	deletingOrder, err := or.GetById(id)
	if err != nil {
		return nil, err
	}
	_, err = ordersCollection.DeleteOne(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	return deletingOrder, nil
}

func (or *OrdersRepository) Update(id int, order model.OrderRequestBody) (*model.Order, error) {
	totalCost := 0
	for _, product := range order.Products {
		totalCost += product.SalePrice
	}
	ordersCollection := or.db.Collection(or.collectionName)
	filter := bson.M{"_id": id}
	update := bson.M{
		"$set": bson.M{"products": order.Products, "lastModifiedDate": time.Now().UTC(), "totalCost": totalCost},
	}
	_, err := ordersCollection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return nil, err
	}
	return or.GetById(id)
}

func (or *OrdersRepository) WriteLog(result any, operation, status string) (int64, error) {

	id, err := or.redis.Incr(context.TODO(), "logOrder:id").Result()
	if err != nil {
		return 0, fmt.Errorf("error incrementing ID: %w", err)
	}

	resultJSON, err := json.Marshal(result)
	if err != nil {
		return 0, fmt.Errorf("error marshaling result: %w", err)
	}

	key := fmt.Sprintf("logOrder:%d", id)
	_, err = or.redis.HSet(context.TODO(), key,
		"id", id,
		"operation", operation,
		"status", status,
		"result", resultJSON,
		"date", time.Now().UTC(),
	).Result()
	if err != nil {
		return 0, fmt.Errorf("error saving log: %w", err)
	}
	_, err = or.redis.Expire(context.TODO(), key, time.Hour*24).Result()
	if err != nil {
		return 0, fmt.Errorf("error setting TTL: %w", err)
	}

	return id, nil
}
