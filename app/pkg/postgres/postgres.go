package postgres

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Config struct {
	Host     string
	Port     string
	Database string
	Username string
	Password string
}

func NewClient(cfg Config) (db *sqlx.DB, err error) {
	connString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", cfg.Username, cfg.Password, cfg.Host,
		cfg.Port, cfg.Database)

	db, err = sqlx.Connect("postgres", connString)
	if err != nil {
		return
	}

	return
}
