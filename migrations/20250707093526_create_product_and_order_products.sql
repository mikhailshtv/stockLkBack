-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';

CREATE SCHEMA IF NOT EXISTS products;

CREATE TABLE IF NOT EXISTS products.products (
    id SERIAL PRIMARY KEY,
    code INTEGER NOT NULL UNIQUE,
    name VARCHAR(255) NOT NULL,
    quantity INTEGER NOT NULL DEFAULT 0 CHECK (quantity >= 0),
    purchase_price INTEGER NOT NULL CHECK (purchase_price >= 0),
    sell_price INTEGER NOT NULL CHECK (sell_price >= 0)
);

CREATE TABLE IF NOT EXISTS orders.order_products (
    order_id INTEGER NOT NULL REFERENCES orders.orders(id) ON DELETE CASCADE,
    product_id INTEGER NOT NULL REFERENCES products.products(id) ON DELETE RESTRICT,
    
    PRIMARY KEY (order_id, product_id)
);

CREATE INDEX IF NOT EXISTS idx_order_products_order ON orders.order_products(order_id);
CREATE INDEX IF NOT EXISTS idx_order_products_product ON orders.order_products(product_id);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';

DROP SCHEMA IF EXISTS products;
DROP TABLE IF EXISTS products.products;
DROP TABLE IF EXISTS orders.order_products;
-- +goose StatementEnd
