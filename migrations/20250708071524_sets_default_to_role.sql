-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';

ALTER TABLE users.users 
ALTER COLUMN role SET DEFAULT 'client';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';

ALTER TABLE users.users 
ALTER COLUMN role DROP DEFAULT;
-- +goose StatementEnd
