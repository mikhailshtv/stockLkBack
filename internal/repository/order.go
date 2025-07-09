package repository

import (
	"golang/stockLkBack/internal/model"

	"github.com/go-redis/redis/v8"
	"github.com/jmoiron/sqlx"
)

type OrdersRepository struct {
	db             *sqlx.DB
	redis          *redis.Client
	collectionName string
}

func NewOrdersRepository(db *sqlx.DB, redis *redis.Client, collectionName string) *OrdersRepository {
	return &OrdersRepository{db: db, redis: redis, collectionName: collectionName}
}

func (or *OrdersRepository) Create(orderRequest model.OrderRequestBody) (*model.Order, error) {
	var order model.Order
	// ordersCollection := or.db.Collection(or.collectionName)
	// _, err := ordersCollection.Find(context.TODO(), bson.D{})
	// if err != nil {
	// 	fmt.Println("Коллекция пуста")
	// 	order.ID = 1
	// 	order.Number = 1
	// } else {
	// 	orderNextID, err := getNextSequence(or.db, "orderid")
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	orderNextNumber, err := getNextSequence(or.db, "orderNumber")
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	order.ID = orderNextID
	// 	order.Number = orderNextNumber
	// }

	// order.CreatedDate = time.Now().UTC()
	// order.LastModifiedDate = time.Now().UTC()
	// order.Status = model.Active
	// var totalCost int32
	// for _, product := range orderRequest.Products {
	// 	totalCost += product.SalePrice
	// }
	// order.TotalCost = totalCost
	// order.Products = orderRequest.Products
	// _, err = ordersCollection.InsertOne(context.TODO(), order)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	return &order, nil
}

func (or *OrdersRepository) GetAll() ([]model.Order, error) {
	// cursor, err := or.db.Collection(or.collectionName).Find(context.TODO(), bson.D{})
	// if err != nil {
	// 	return nil, err
	// }
	var orders []model.Order
	// if err = cursor.All(context.TODO(), &orders); err != nil {
	// 	return nil, err
	// }
	return orders, nil
}

func (or *OrdersRepository) GetByID(id int32) (*model.Order, error) {
	// ordersCollection := or.db.Collection(or.collectionName)
	// var result bson.M
	// filter := bson.M{"_id": id}
	// err := ordersCollection.FindOne(context.TODO(), filter).Decode(&result)
	// if err != nil {
	// 	return nil, errors.New(NotFoundErrorMessage)
	// }
	var order model.Order
	// resultBytes, err := bson.Marshal(result)
	// if err != nil {
	// 	return nil, err
	// }
	// err = bson.Unmarshal(resultBytes, &order)
	// if err != nil {
	// 	return nil, err
	// }

	return &order, nil
}

func (or *OrdersRepository) Delete(id int32) (*model.Order, error) {
	// ordersCollection := or.db.Collection(or.collectionName)
	// filter := bson.M{"_id": id}
	deletingOrder, _ := or.GetByID(id)
	// if err != nil {
	// 	return nil, err
	// }
	// _, err = ordersCollection.DeleteOne(context.TODO(), filter)
	// if err != nil {
	// 	return nil, err
	// }
	return deletingOrder, nil
}

func (or *OrdersRepository) Update(id int32, order model.OrderRequestBody) (*model.Order, error) {
	var totalCost int32
	for _, product := range order.Products {
		totalCost += product.SellPrice
	}
	// ordersCollection := or.db.Collection(or.collectionName)
	// filter := bson.M{"_id": id}
	// update := bson.M{
	// 	"$set": bson.M{"products": order.Products, "lastModifiedDate": time.Now().UTC(), "totalCost": totalCost},
	// }
	// _, err := ordersCollection.UpdateOne(context.TODO(), filter, update)
	// if err != nil {
	// 	return nil, err
	// }
	return or.GetByID(id)
}

func (or *OrdersRepository) WriteLog(result any, operation, status, tableName string) (int64, error) {
	return WriteLog(result, operation, status, tableName, or.redis)
}
