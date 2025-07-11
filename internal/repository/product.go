package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/mikhailshtv/stockLkBack/internal/model"

	"github.com/go-redis/redis/v8"
	"github.com/jmoiron/sqlx"
)

type ProductsRepository struct {
	db    *sqlx.DB
	redis *redis.Client
}

func NewProductsRepository(db *sqlx.DB, redis *redis.Client) *ProductsRepository {
	return &ProductsRepository{db: db, redis: redis}
}

func (pr *ProductsRepository) Create(ctx context.Context, product model.Product) (*model.Product, error) {
	const query = `
		INSERT INTO products.products (
			code,
			name,
			quantity,
			purchase_price,
			sell_price
		) VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`
	err := pr.db.QueryRowContext(
		ctx,
		query,
		product.Code,
		product.Name,
		product.Quantity,
		product.PurchasePrice,
		product.SellPrice,
	).Scan(&product.ID)
	if err != nil {
		if isDuplicateKeyError(err) {
			return nil, fmt.Errorf("продукт с кодом %d уже существует", product.Code)
		}
		return nil, fmt.Errorf("ошибка при создании продукта: %w", err)
	}

	return &product, nil
}

func (pr *ProductsRepository) GetAll(ctx context.Context) ([]model.Product, error) {
	const query = `SELECT * FROM products.products`

	var products []model.Product
	err := pr.db.SelectContext(ctx, &products, query)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return []model.Product{}, nil
		}
		return nil, fmt.Errorf("ошибка при получении списка продуктов: %w", err)
	}

	return products, nil
}

func (pr *ProductsRepository) GetByID(ctx context.Context, id int32) (*model.Product, error) {
	var product model.Product
	err := pr.db.GetContext(ctx, &product,
		"SELECT * FROM products.products WHERE id = $1", id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("продукт не найден: %w", err)
		}
		return nil, fmt.Errorf("ошибка при получении продукта: %w", err)
	}
	return &product, nil
}

func (pr *ProductsRepository) Delete(ctx context.Context, id int32) (*model.Product, error) {
	const query = `
		WITH deleted AS (
			DELETE FROM products.products 
			WHERE id = $1
			RETURNING *
		)
		SELECT * FROM deleted
	`

	deletedProduct := model.Product{}
	err := pr.db.QueryRowxContext(ctx, query, id).StructScan(&deletedProduct)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("продукт не найден: %w", err)
		}
		return nil, fmt.Errorf("ошибка удаления продукта: %w", err)
	}

	return &deletedProduct, nil
}

func (pr *ProductsRepository) Update(ctx context.Context, id int32, product model.Product) (*model.Product, error) {
	const query = `
		UPDATE products.products SET
			code = $1,
			name = $2,
			quantity = $3,
			purchase_price = $4,
			sell_price = $5
		WHERE id = $6
		RETURNING *
	`

	updatedProduct := model.Product{}
	err := pr.db.QueryRowxContext(
		ctx,
		query,
		product.Code,
		product.Name,
		product.Quantity,
		product.PurchasePrice,
		product.SellPrice,
		id,
	).StructScan(&updatedProduct)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("продукт не найден: %w", err)
		}
		if isDuplicateKeyError(err) {
			return nil, fmt.Errorf("продукт с кодом %d уже существует", product.Code)
		}
		return nil, fmt.Errorf("ошибка обновления продукта: %w", err)
	}

	return &updatedProduct, nil
}

func (pr *ProductsRepository) WriteLog(result any, operation, status, tableName string) (int64, error) {
	return WriteLog(result, operation, status, tableName, pr.redis)
}
