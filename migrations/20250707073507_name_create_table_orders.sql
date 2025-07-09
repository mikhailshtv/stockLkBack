-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';

CREATE TABLE IF NOT EXISTS orders.orders(
  id SERIAL PRIMARY KEY,
  order_number INTEGER NOT NULL,
  total_cost INTEGER NOT NULL,
  created_date TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  last_modified_date TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  status VARCHAR(20) NOT NULL CHECK (status IN ('active', 'executed', 'deleted'))
);

COMMENT ON TABLE orders.orders IS 'Таблица для хранения заказов';
COMMENT ON COLUMN orders.orders.status IS 'Статус заказа: active, executed, deleted';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';

DROP TABLE IF EXISTS orders.orders
-- +goose StatementEnd
