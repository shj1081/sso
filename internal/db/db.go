package db

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

// NewDatabase는 DSN을 받아 MySQL 연결을 생성합니다.
func NewDatabase(dsn string) (*sqlx.DB, error) {
	db, err := sqlx.Connect("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("error connecting to MySQL: %v", err)
	}
	return db, nil
}
