package postgres

import (
	"github.com/uptrace/bun"
)

type DB struct {
	db *bun.DB
}

func NewDB(db *bun.DB) *DB {
	return &DB{db: db}
}
