package postgres

import (
	"context"
	"fmt"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"
)

type Config struct {
	Host     string
	Port     string
	Database string
	Username string
	Password string
}

type Client interface {
	Exec(ctx context.Context, query string, args ...any) (pgconn.CommandTag, error)
	Get(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	Select(ctx context.Context, dest interface{}, query string, args ...interface{}) error
}

type client struct {
	pool *pgxpool.Pool
}

func NewClient(ctx context.Context, cfg Config) (Client, error) {
	config, err := pgxpool.ParseConfig(fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", cfg.Username, cfg.Password, cfg.Host,
		cfg.Port, cfg.Database))
	if err != nil {
		return nil, err
	}

	var pool *pgxpool.Pool
	pool, err = pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, err
	}

	return &client{
		pool: pool,
	}, nil
}

func (c *client) Exec(ctx context.Context, query string, args ...any) (pgconn.CommandTag, error) {
	return c.pool.Exec(ctx, query, args...)
}

func (c *client) Get(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	return pgxscan.Get(ctx, c.pool, dest, query, args...)
}

func (c *client) Select(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	return pgxscan.Select(ctx, c.pool, dest, query, args...)
}
