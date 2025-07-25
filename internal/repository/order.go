package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/mikhailshtv/stockLkBack/internal/model"

	"github.com/go-redis/redis/v8"
	"github.com/jackc/pgconn"
	"github.com/jmoiron/sqlx"
)

const (
	maxRetries = 3                      // Максимальное число попыток
	retryDelay = 100 * time.Millisecond // Задержка между попытками
)

type OrdersRepository struct {
	db             *sqlx.DB
	redis          *redis.Client
	collectionName string
}

func NewOrdersRepository(db *sqlx.DB, redis *redis.Client, collectionName string) *OrdersRepository {
	return &OrdersRepository{db: db, redis: redis, collectionName: collectionName}
}

func (or *OrdersRepository) Create(
	ctx context.Context,
	orderRequest model.OrderRequestBody,
	userID int,
) (*model.Order, error) {
	if len(orderRequest.Products) == 0 {
		return nil, fmt.Errorf("список товаров не может быть пустым")
	}
	var lastErr error

	for i := 0; i < maxRetries; i++ {
		order, err := or.tryCreateOrder(ctx, orderRequest, userID)
		if err == nil {
			return order, nil
		}

		// Проверяем, нужно ли повторять (ошибка сериализации)
		lastErr = err
		if !isRetryableError(err) {
			return nil, err
		}

		time.Sleep(retryDelay)
	}

	return nil, fmt.Errorf("не удалось создать заказ после %d попыток: %w", maxRetries, lastErr)
}

func (or *OrdersRepository) GetAll(ctx context.Context, userID int, role model.UserRole) ([]model.Order, error) {
	var orders []model.Order

	query := `
		SELECT *
		FROM orders.orders o
	`

	var builder strings.Builder
	builder.WriteString(query)
	args := []interface{}{}

	if role != model.RoleEmployee {
		builder.WriteString(" WHERE o.user_id = $1")
		args = append(args, userID)
	}

	builder.WriteString(" ORDER BY id ASC")
	query = builder.String()
	// Получаем заказы
	err := or.db.SelectContext(ctx, &orders, query, args...)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения списка заказов: %w", err)
	}

	// Для каждого заказа получаем список товаров
	for i := range orders {
		err = or.db.SelectContext(ctx, &orders[i].Products, `
			SELECT
				p.id,
				p.code,
				p.name,
				op.quantity,
				op.sell_price
			FROM orders.order_products op
			JOIN products.products p ON op.product_id = p.id
			WHERE op.order_id = $1
		`, orders[i].ID)
		if err != nil {
			return nil, fmt.Errorf("ошибка получения товаров для заказа %d: %w", orders[i].ID, err)
		}
	}

	return orders, nil
}

func (or *OrdersRepository) GetByID(ctx context.Context, id, userID int, role model.UserRole) (*model.Order, error) {
	query := `
		SELECT *
		FROM orders.orders
		WHERE id = $1
	`
	args := []interface{}{id}

	var builder strings.Builder
	if role != model.RoleEmployee {
		builder.WriteString(query)
		builder.WriteString(" AND user_id = $2")
		args = append(args, userID)
		query = builder.String()
	}

	// Получаем заказ
	var order model.Order
	err := or.db.GetContext(ctx, &order, query, args...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("заказ не найден или не принадлежит пользователю")
		}
		return nil, fmt.Errorf("ошибка получения заказа: %w", err)
	}

	// Получаем товары для заказа
	err = or.db.SelectContext(ctx, &order.Products, `
		SELECT 
			p.id,
			p.code,
			p.name,
			op.quantity,
			op.sell_price
		FROM orders.order_products op
		JOIN products.products p ON op.product_id = p.id
		WHERE op.order_id = $1
	`, id)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения товаров заказа: %w", err)
	}

	return &order, nil
}

func (or *OrdersRepository) Delete(ctx context.Context, id, userID int) (*model.Order, error) {
	var lastErr error

	for i := 0; i < maxRetries; i++ {
		order, err := or.tryDeleteOrder(ctx, id, userID)
		if err == nil {
			return order, nil
		}

		lastErr = err
		if !isRetryableError(err) {
			return nil, err
		}

		time.Sleep(retryDelay)
	}

	return nil, fmt.Errorf("не удалось удалить заказ после %d попыток: %w", maxRetries, lastErr)
}

func (or *OrdersRepository) Update(
	ctx context.Context,
	id int,
	orderRequest model.OrderRequestBody,
	userID int,
) (*model.Order, error) {
	if len(orderRequest.Products) == 0 {
		return nil, fmt.Errorf("список товаров не может быть пустым")
	}

	var lastErr error

	for i := 0; i < maxRetries; i++ {
		order, err := or.tryUpdateOrder(ctx, id, orderRequest, userID)
		if err == nil {
			return order, nil
		}

		lastErr = err
		if !isRetryableError(err) {
			return nil, err
		}

		time.Sleep(retryDelay)
	}

	return nil, fmt.Errorf("не удалось обновить заказ после %d попыток: %w", maxRetries, lastErr)
}

func (or *OrdersRepository) UpdateStatus(
	ctx context.Context,
	id int,
	orderStatusRequest model.OrderStatusRequest,
	userID int,
) (*model.Order, error) {
	var lastErr error

	for i := 0; i < maxRetries; i++ {
		order, err := or.tryUpdateStatus(ctx, id, orderStatusRequest, userID)
		if err == nil {
			return order, nil
		}

		lastErr = err
		if !isRetryableError(err) {
			return nil, err
		}

		time.Sleep(retryDelay)
	}

	return nil, fmt.Errorf("не удалось обновить статус заказа после %d попыток: %w", maxRetries, lastErr)
}

func (or *OrdersRepository) WriteLog(result any, operation, status, tableName string) (int64, error) {
	return WriteLog(result, operation, status, tableName, or.redis)
}

func (or *OrdersRepository) tryCreateOrder(
	ctx context.Context,
	request model.OrderRequestBody,
	userID int,
) (*model.Order, error) {
	tx, err := or.db.BeginTxx(ctx, &sql.TxOptions{
		Isolation: sql.LevelSerializable,
	})
	if err != nil {
		return nil, fmt.Errorf("ошибка начала транзакции: %w", err)
	}
	defer tx.Rollback()

	var order model.Order

	// 1. Создаем запись заказа
	err = tx.QueryRowxContext(ctx, `
		INSERT INTO orders.orders (
			user_id, 
			order_number, 
			status,
			total_cost
		) VALUES (
			$1, 
			(SELECT COALESCE(MAX(order_number), 0) + 1 FROM orders.orders),
			'active',
			0
		)
		RETURNING id, order_number, status, created_date, last_modified_date, user_id
	`, userID).StructScan(&order)
	if err != nil {
		return nil, fmt.Errorf("ошибка создания заказа: %w", err)
	}

	// 2. Добавляем товары в заказ
	for _, product := range request.Products {
		log.Println(product.ProductID)
		var available int
		var version int
		err = tx.QueryRowContext(ctx, `
			SELECT quantity, version FROM products.products WHERE id = $1
		`, product.ProductID).Scan(&available, &version)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, fmt.Errorf("товар с ID %d не найден", product.ProductID)
			}
			return nil, fmt.Errorf("ошибка проверки товара с ID %d: %w", product.ProductID, err)
		}

		if available < product.Quantity {
			return nil, fmt.Errorf("недостаточно товара с ID %d (доступно: %d)", product.ProductID, available)
		}

		res, err := tx.ExecContext(ctx, `
			UPDATE products.products 
			SET quantity = quantity - $1, version = version + 1 
			WHERE id = $2 AND version = $3
		`, product.Quantity, product.ProductID, version)
		if err != nil {
			return nil, fmt.Errorf("ошибка обновления остатков: %w", err)
		}

		if rowsAffected, _ := res.RowsAffected(); rowsAffected == 0 {
			return nil, fmt.Errorf("конфликт версий товара %d (параллельное изменение)", product.ProductID)
		}

		_, err = tx.ExecContext(ctx, `
			INSERT INTO orders.order_products (
				order_id, 
				product_id, 
				quantity, 
				sell_price
			) VALUES ($1, $2, $3, $4)
		`, order.ID, product.ProductID, product.Quantity, product.SellPrice)
		if err != nil {
			return nil, fmt.Errorf("ошибка добавления товара в заказ: %w", err)
		}
	}

	// 3. Получаем полные данные заказа
	err = tx.GetContext(ctx, &order, `
		SELECT 
			id, 
			order_number, 
			total_cost, 
			created_date, 
			last_modified_date, 
			status,
			user_id
		FROM orders.orders 
		WHERE id = $1
	`, order.ID)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения заказа: %w", err)
	}

	// 4. Получаем товары заказа
	var products []model.Product
	err = tx.SelectContext(ctx, &products, `
		SELECT 
			p.id,
			p.code,
			p.name,
			op.quantity,
			op.sell_price
		FROM orders.order_products op
		JOIN products.products p ON op.product_id = p.id
		WHERE op.order_id = $1
	`, order.ID)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения товаров заказа: %w", err)
	}

	order.Products = products

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("ошибка фиксации транзакции: %w", err)
	}

	return &order, nil
}

func (or *OrdersRepository) tryUpdateOrder(
	ctx context.Context,
	id int,
	orderRequest model.OrderRequestBody,
	userID int,
) (*model.Order, error) {
	tx, err := or.db.BeginTxx(ctx, &sql.TxOptions{
		Isolation: sql.LevelSerializable,
	})
	if err != nil {
		return nil, fmt.Errorf("ошибка начала транзакции: %w", err)
	}
	defer tx.Rollback()

	order, err := or.getOrderForUpdate(ctx, tx, id, userID)
	if err != nil {
		return nil, err
	}

	currentProducts, err := or.getCurrentOrderProducts(ctx, tx, order.ID)
	if err != nil {
		return nil, err
	}

	err = or.processProductChanges(ctx, tx, order.ID, currentProducts, orderRequest.Products)
	if err != nil {
		return nil, err
	}

	err = or.addNewProducts(ctx, tx, order.ID, currentProducts, orderRequest.Products)
	if err != nil {
		return nil, err
	}

	err = or.updateOrderModifiedDate(ctx, tx, order.ID)
	if err != nil {
		return nil, err
	}

	updatedOrder, err := or.getUpdatedOrder(ctx, tx, order.ID)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("ошибка фиксации транзакции: %w", err)
	}

	return updatedOrder, nil
}

func (or *OrdersRepository) getOrderForUpdate(
	ctx context.Context,
	tx *sqlx.Tx,
	id, userID int,
) (*model.Order, error) {
	var order model.Order
	err := tx.GetContext(ctx, &order, `
		SELECT *
		FROM orders.orders 
		WHERE id = $1 AND user_id = $2
		FOR UPDATE
	`, id, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("заказ не найден или не принадлежит пользователю")
		}
		return nil, fmt.Errorf("ошибка получения заказа: %w", err)
	}
	return &order, nil
}

func (or *OrdersRepository) getCurrentOrderProducts(
	ctx context.Context,
	tx *sqlx.Tx,
	orderID int,
) ([]model.OrderProduct, error) {
	var products []model.OrderProduct
	err := tx.SelectContext(ctx, &products, `
		SELECT product_id, quantity, sell_price
		FROM orders.order_products
		WHERE order_id = $1
	`, orderID)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения текущих товаров: %w", err)
	}
	return products, nil
}

func (or *OrdersRepository) processProductChanges(
	ctx context.Context,
	tx *sqlx.Tx,
	orderID int,
	currentProducts []model.OrderProduct,
	newProducts []model.OrderProduct,
) error {
	oldProductsMap := make(map[int]model.OrderProduct)
	for _, p := range currentProducts {
		oldProductsMap[p.ProductID] = p
	}

	newProductsMap := make(map[int]model.OrderProduct)
	for _, p := range newProducts {
		newProductsMap[p.ProductID] = p
	}

	// 4. Обрабатываем изменения товаров
	for productID, oldProduct := range oldProductsMap {
		newProduct, exists := newProductsMap[productID]

		// Товар удален из заказа - возвращаем остатки
		if !exists {
			_, err := tx.ExecContext(ctx, `
				UPDATE products.products
				SET quantity = quantity + $1
				WHERE id = $2
			`, oldProduct.Quantity, productID)
			if err != nil {
				return fmt.Errorf("ошибка возврата товара %d: %w", productID, err)
			}

			// Удаляем товар из заказа
			_, err = tx.ExecContext(ctx, `
				DELETE FROM orders.order_products
				WHERE order_id = $1 AND product_id = $2
			`, orderID, productID)
			if err != nil {
				return fmt.Errorf("ошибка удаления товара %d: %w", productID, err)
			}
			continue
		}

		// Количество изменилось - корректируем остатки
		if oldProduct.Quantity != newProduct.Quantity {
			diff := oldProduct.Quantity - newProduct.Quantity
			_, err := tx.ExecContext(ctx, `
				UPDATE products.products
				SET quantity = quantity + $1
				WHERE id = $2
			`, diff, productID)
			if err != nil {
				return fmt.Errorf("ошибка обновления количества товара %d: %w", productID, err)
			}

			// Обновляем количество в заказе
			_, err = tx.ExecContext(ctx, `
				UPDATE orders.order_products
				SET quantity = $1, sell_price = $2
				WHERE order_id = $3 AND product_id = $4
			`, newProduct.Quantity, newProduct.SellPrice, orderID, productID)
			if err != nil {
				return fmt.Errorf("ошибка обновления товара %d в заказе: %w", productID, err)
			}
		}
	}
	return nil
}

func (or *OrdersRepository) addNewProducts(
	ctx context.Context,
	tx *sqlx.Tx,
	orderID int,
	currentProducts []model.OrderProduct,
	newProducts []model.OrderProduct,
) error {
	oldProductsMap := make(map[int]model.OrderProduct)
	for _, p := range currentProducts {
		oldProductsMap[p.ProductID] = p
	}

	for _, newProduct := range newProducts {
		if _, exists := oldProductsMap[newProduct.ProductID]; !exists {
			// Проверяем доступность товара
			var available int
			err := tx.GetContext(ctx, &available, `
				SELECT quantity FROM products.products WHERE id = $1
			`, newProduct.ProductID)
			if err != nil {
				return fmt.Errorf("ошибка проверки товара %d: %w", newProduct.ProductID, err)
			}

			if available < newProduct.Quantity {
				return fmt.Errorf("недостаточно товара %d (доступно: %d, требуется: %d)",
					newProduct.ProductID, available, newProduct.Quantity)
			}

			// Резервируем товар
			_, err = tx.ExecContext(ctx, `
				UPDATE products.products
				SET quantity = quantity - $1
				WHERE id = $2
			`, newProduct.Quantity, newProduct.ProductID)
			if err != nil {
				return fmt.Errorf("ошибка резервирования товара %d: %w", newProduct.ProductID, err)
			}

			// Добавляем в заказ
			_, err = tx.ExecContext(ctx, `
				INSERT INTO orders.order_products
				(order_id, product_id, quantity, sell_price)
				VALUES ($1, $2, $3, $4)
			`, orderID, newProduct.ProductID, newProduct.Quantity, newProduct.SellPrice)
			if err != nil {
				return fmt.Errorf("ошибка добавления товара %d: %w", newProduct.ProductID, err)
			}
		}
	}

	return nil
}

func (or *OrdersRepository) updateOrderModifiedDate(ctx context.Context, tx *sqlx.Tx, orderID int) error {
	_, err := tx.ExecContext(ctx, `
		UPDATE orders.orders
		SET last_modified_date = NOW()
		WHERE id = $1
	`, orderID)
	if err != nil {
		return fmt.Errorf("ошибка обновления даты заказа: %w", err)
	}
	return nil
}

func (or *OrdersRepository) getUpdatedOrder(ctx context.Context, tx *sqlx.Tx, orderID int) (*model.Order, error) {
	var order model.Order
	err := tx.GetContext(ctx, &order, `
		SELECT *
		FROM orders.orders 
		WHERE id = $1
	`, orderID)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения обновленного заказа: %w", err)
	}

	var products []model.Product
	err = tx.SelectContext(ctx, &products, `
		SELECT p.id, p.code, p.name, op.quantity, op.sell_price
		FROM orders.order_products op
		JOIN products.products p ON op.product_id = p.id
		WHERE op.order_id = $1
	`, orderID)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения товаров заказа: %w", err)
	}

	order.Products = products
	return &order, nil
}

func (or *OrdersRepository) tryDeleteOrder(ctx context.Context, orderID, userID int) (*model.Order, error) {
	tx, err := or.db.BeginTxx(ctx, &sql.TxOptions{
		Isolation: sql.LevelSerializable,
	})
	if err != nil {
		return nil, fmt.Errorf("ошибка начала транзакции: %w", err)
	}
	defer tx.Rollback()

	// 1. Получаем заказ
	var order model.Order
	err = tx.GetContext(ctx, &order, `
		SELECT *
		FROM orders.orders 
		WHERE id = $1 AND user_id = $2
		FOR UPDATE
	`, orderID, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("заказ не найден или не принадлежит пользователю")
		}
		return nil, fmt.Errorf("ошибка получения заказа: %w", err)
	}

	// 2. Если заказ уже удален, просто возвращаем его
	if order.Status.Key == "deleted" {
		return &order, nil
	}

	// 3. Возвращаем товары на склад (если заказ активен)
	if order.Status.Key == "active" {
		// Получаем все товары из заказа
		var products []model.OrderProduct
		err = tx.SelectContext(ctx, &products, `
			SELECT product_id, quantity
			FROM orders.order_products
			WHERE order_id = $1
		`, order.ID)
		if err != nil {
			return nil, fmt.Errorf("ошибка получения товаров заказа: %w", err)
		}

		// Возвращаем каждый товар на склад
		for _, product := range products {
			_, err = tx.ExecContext(ctx, `
				UPDATE products.products
				SET quantity = quantity + $1
				WHERE id = $2
			`, product.Quantity, product.ProductID)
			if err != nil {
				return nil, fmt.Errorf("ошибка возврата товара %d: %w", product.ProductID, err)
			}
		}
	}

	// 4. Помечаем заказ как удаленный
	_, err = tx.ExecContext(ctx, `
		UPDATE orders.orders
		SET status = 'deleted',
			last_modified_date = NOW()
		WHERE id = $1
	`, order.ID)
	if err != nil {
		return nil, fmt.Errorf("ошибка обновления статуса заказа: %w", err)
	}

	// 5. Получаем обновленный заказ
	err = tx.GetContext(ctx, &order, `
		SELECT *
		FROM orders.orders 
		WHERE id = $1
	`, order.ID)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения обновленного заказа: %w", err)
	}

	// 6. Получаем товары заказа (для возврата в ответе и последующего логирования)
	var products []model.Product
	err = tx.SelectContext(ctx, &products, `
		SELECT p.id, p.code, p.name, op.quantity, op.sell_price
		FROM orders.order_products op
		JOIN products.products p ON op.product_id = p.id
		WHERE op.order_id = $1
	`, order.ID)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения товаров заказа: %w", err)
	}

	order.Products = products

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("ошибка фиксации транзакции: %w", err)
	}

	return &order, nil
}

func (or *OrdersRepository) tryUpdateStatus(
	ctx context.Context,
	id int,
	orderStatusRequest model.OrderStatusRequest,
	userID int,
) (*model.Order, error) {
	tx, err := or.db.BeginTxx(ctx, &sql.TxOptions{
		Isolation: sql.LevelSerializable,
	})
	if err != nil {
		return nil, fmt.Errorf("ошибка начала транзакции: %w", err)
	}
	defer tx.Rollback()

	// 1. Получаем текущий заказ с блокировкой
	var order model.Order
	query := `
		SELECT id, order_number, status, user_id 
		FROM orders.orders 
		WHERE id = $1
		AND user_id = $2
		FOR UPDATE
	`

	err = tx.GetContext(ctx, &order, query, id, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("заказ не найден или доступ запрещен")
		}
		return nil, fmt.Errorf("ошибка получения заказа: %w", err)
	}

	// 2. Проверяем допустимость изменения статуса
	if !isValidStatusTransition(order.Status, orderStatusRequest.Status) {
		return nil, fmt.Errorf("недопустимый переход статуса из %s в %s",
			order.Status.Key, orderStatusRequest.Status.Key)
	}

	// 3. Обновляем статус
	_, err = tx.ExecContext(ctx, `
		UPDATE orders.orders 
		SET status = $1, last_modified_date = NOW()
		WHERE id = $2
	`, orderStatusRequest.Status, id)
	if err != nil {
		return nil, fmt.Errorf("ошибка обновления статуса: %w", err)
	}

	// 4. Получаем обновленный заказ с товарами
	err = tx.GetContext(ctx, &order, `
		SELECT id, order_number, status, total_cost, 
			created_date, last_modified_date, user_id
		FROM orders.orders 
		WHERE id = $1
	`, id)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения заказа: %w", err)
	}

	err = tx.SelectContext(ctx, &order.Products, `
		SELECT p.id, p.code, p.name, op.quantity, op.sell_price
		FROM orders.order_products op
		JOIN products.products p ON op.product_id = p.id
		WHERE op.order_id = $1
	`, id)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения товаров: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("ошибка фиксации транзакции: %w", err)
	}

	return &order, nil
}

// Пока единственный допустимый переход active -> executed,
// но так хотя бы можно расширить варианты переходов, если добавятся новые статусы.
func isValidStatusTransition(oldStatus, newStatus model.OrderStatus) bool {
	validTransitions := map[string][]string{
		"active":   {"executed"},
		"executed": {},
		"deleted":  {},
	}

	allowed, exists := validTransitions[oldStatus.Key]
	if !exists {
		return false
	}

	for _, s := range allowed {
		if s == newStatus.Key {
			return true
		}
	}
	return false
}

func isRetryableError(err error) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code == "40001"
	}
	return false
}
