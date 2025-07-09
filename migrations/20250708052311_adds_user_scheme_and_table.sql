-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';

CREATE SCHEMA IF NOT EXISTS users;

-- Создаем таблицу users
CREATE TABLE IF NOT EXISTS users.users (
    id SERIAL PRIMARY KEY,
    login VARCHAR(50) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    role VARCHAR(20) NOT NULL CHECK (role IN ('client', 'employee'))
);

CREATE INDEX IF NOT EXISTS idx_users_email ON users.users(email);

COMMENT ON SCHEMA users IS 'Схема для хранения данных пользователей';
COMMENT ON TABLE users.users IS 'Таблица пользователей системы';
COMMENT ON COLUMN users.users.role IS 'Роль пользователя: client или employee';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';

DROP SCHEMA IF EXISTS users;
DROP TABLE IF EXISTS users.users;
-- +goose StatementEnd
