-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';

CREATE INDEX IF NOT EXISTS idx_orders_user ON orders.orders(user_id);

CREATE OR REPLACE FUNCTION update_order_total()
RETURNS TRIGGER AS $$
BEGIN
    UPDATE orders.orders
    SET total_cost = (
        SELECT SUM(op.quantity * op.sell_price)
        FROM orders.order_products op
        WHERE op.order_id = NEW.order_id
    )
    WHERE id = NEW.order_id;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_order_products_update
AFTER INSERT OR UPDATE OR DELETE ON orders.order_products
FOR EACH ROW EXECUTE FUNCTION update_order_total();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';

DROP INDEX IF EXISTS idx_orders_user ON orders.orders(user_id);
DROP TRIGGER IF EXISTS trg_order_products_update ON orders.order_products;
DROP FUNCTION IF EXISTS update_order_total;
-- +goose StatementEnd
