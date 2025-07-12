-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';

CREATE OR REPLACE FUNCTION check_product_delete()
RETURNS TRIGGER AS $$
BEGIN
    IF EXISTS (
        SELECT 1 FROM orders.order_products op
        JOIN orders.orders o ON op.order_id = o.id
        WHERE op.product_id = OLD.id 
        AND o.status = 'active'
    ) THEN
        RAISE EXCEPTION 'Cannot delete product: it is used in active orders';
    END IF;
    RETURN OLD;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER prevent_active_product_delete
BEFORE DELETE ON products.products
FOR EACH ROW EXECUTE FUNCTION check_product_delete();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';

DROP TRIGGER IF EXISTS prevent_active_product_delete ON products.products;
DROP FUNCTION IF EXISTS check_product_delete;
-- +goose StatementEnd
