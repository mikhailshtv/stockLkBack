-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';

ALTER TABLE products.products ADD COLUMN version INTEGER NOT NULL DEFAULT 0;
COMMENT ON COLUMN products.products.version IS 'Версия для optimistic lock';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';

ALTER TABLE products.products DROP COLUMN version;
-- +goose StatementEnd
