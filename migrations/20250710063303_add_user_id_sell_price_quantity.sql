-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';

ALTER TABLE orders.orders
ADD COLUMN user_id INTEGER NOT NULL,
ADD CONSTRAINT fk_orders_user
    FOREIGN KEY (user_id) 
    REFERENCES users.users(id)
    ON DELETE RESTRICT;

ALTER TABLE orders.order_products
ADD COLUMN quantity INTEGER NOT NULL DEFAULT 1 CHECK (quantity > 0),
ADD COLUMN sell_price INTEGER NOT NULL CHECK (sell_price >= 0);

COMMENT ON COLUMN orders.order_products.quantity IS 'Количество товара в заказе';
COMMENT ON COLUMN orders.order_products.sell_price IS 'Цена товара на момент создания заказа';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';

ALTER TABLE orders.orders
DROP CONSTRAINT fk_orders_user,
DROP COLUMN user_id;

ALTER TABLE orders.order_products
DROP COLUMN quantity,
DROP COLUMN sell_price;
-- +goose StatementEnd
