package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

type SQLConfig struct {
	Host           string
	Port           string
	User           string
	Password       string
	Database       string
	SSLMode        string
	MaxConnections int
	Timeout        int
}

func (c SQLConfig) CreateDsn() string {
	return fmt.Sprintf(
		"%s://%s:%s@%s:%s/%s?sslmode=%s&pool_max_conns=%d&connect_timeout=%d",
		"postgres",
		c.User,
		c.Password,
		c.Host,
		c.Port,
		c.Database,
		c.SSLMode,
		c.MaxConnections,
		c.Timeout,
	)
}

func NewSqlxConn(ctx context.Context, dsn string) (*sqlx.DB, error) {
	cfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("pgx parse config: %w", err)
	}

	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("connect to database: %w", err)
	}

	sqldb := stdlib.OpenDBFromPool(pool)

	db := sqlx.NewDb(sqldb, "pgxpool")

	return db, nil
}
