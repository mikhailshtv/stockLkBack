package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

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

func (pr *ProductsRepository) GetAll(ctx context.Context, params model.ProductQueryParams) ([]model.Product, error) {
	baseQuery := `SELECT * FROM products.products WHERE 1=1`
	// Строим запрос с фильрами.
	query, args := pr.buildProductsQuery(baseQuery, params)

	// Сортировка.
	validSortFields := map[string]bool{
		"id":             true,
		"code":           true,
		"quantity":       true,
		"name":           true,
		"purchase_price": true,
		"sell_price":     true,
	}

	if params.SortField != "" && validSortFields[params.SortField] {
		if params.SortOrder == "" {
			params.SortOrder = sortAscParam
		}
		params.SortOrder = strings.ToUpper(params.SortOrder)
		if params.SortOrder != sortAscParam && params.SortOrder != sortDescParam {
			params.SortOrder = sortAscParam
		}
	} else {
		params.SortField = "id"
		params.SortOrder = sortAscParam
	}
	query += fmt.Sprintf(" ORDER BY %s %s", params.SortField, params.SortOrder)

	// Пагинация.
	if params.PageSize > 0 {
		offset := (params.Page - 1) * params.PageSize
		query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", len(args)+1, len(args)+2)
		args = append(args, params.PageSize, offset)
	}

	var products []model.Product
	err := pr.db.SelectContext(ctx, &products, query, args...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return []model.Product{}, nil
		}
		return nil, fmt.Errorf("ошибка при получении списка продуктов: %w", err)
	}

	return products, nil
}

func (pr *ProductsRepository) GetTotalCount(ctx context.Context, params model.ProductQueryParams) (int, error) {
	baseQuery := `SELECT COUNT(*) FROM products.products WHERE 1=1`
	query, args := pr.buildProductsQuery(baseQuery, params)

	var total int
	err := pr.db.GetContext(ctx, &total, query, args...)
	if err != nil {
		return 0, fmt.Errorf("ошибка при получении общего количества продуктов: %w", err)
	}

	return total, nil
}

func (pr *ProductsRepository) GetByID(ctx context.Context, id int) (*model.Product, error) {
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

func (pr *ProductsRepository) Delete(ctx context.Context, id int) (*model.Product, error) {
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

func (pr *ProductsRepository) Update(ctx context.Context, id int, product model.Product) (*model.Product, error) {
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

func (pr *ProductsRepository) buildProductsQuery(baseQuery string, params model.ProductQueryParams) (string, []any) {
	query := baseQuery
	args := []any{}
	argPos := 1

	if params.Code != nil {
		query += fmt.Sprintf(" AND code = $%d", argPos)
		args = append(args, *params.Code)
		argPos++
	}

	if params.Quantity != nil {
		query += fmt.Sprintf(" AND quantity = $%d", argPos)
		args = append(args, *params.Quantity)
		argPos++
	}

	if params.Name != "" {
		query += fmt.Sprintf(" AND name ILIKE $%d", argPos)
		args = append(args, "%"+params.Name+"%")
		argPos++
	}

	if params.PurchasePrice != nil {
		query += fmt.Sprintf(" AND purchase_price = $%d", argPos)
		args = append(args, *params.PurchasePrice)
		argPos++
	}

	if params.SellPrice != nil {
		query += fmt.Sprintf(" AND sell_price = $%d", argPos)
		args = append(args, *params.SellPrice)
	}

	return query, args
}

func (pr *ProductsRepository) WriteLog(result any, operation, status, tableName string) (int64, error) {
	return WriteLog(result, operation, status, tableName, pr.redis)
}
