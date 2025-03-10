package db

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/shj1081/sso/config"
)

type Database struct {
	db *sqlx.DB
}

func NewDatabase() (*Database, error) {

	cfg := config.LoadConfig()

	db, err := sqlx.Open(cfg.DBDriver, cfg.DBURL)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %w", err)
	}

	return &Database{db: db}, nil
}

func (d *Database) Close() error {
	return d.db.Close()
}

func (d *Database) GetDB() *sqlx.DB {
	return d.db
}
