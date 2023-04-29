package postgres

import (
	"database/sql"
	"fmt"

	"github.com/parwin-pp/todo-application/internal/config"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

type DB struct {
	db *bun.DB
}

func (db *DB) Close() error {
	return db.db.Close()
}

func NewDB(db *bun.DB) *DB {
	return &DB{db: db}
}

func DSN(config config.DatabaseConfig) string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		config.User,
		config.Password,
		config.Host,
		config.Port,
		config.Name,
	)
}

func GetConnection(config config.DatabaseConfig) *bun.DB {
	dsn := DSN(config)
	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))
	return bun.NewDB(sqldb, pgdialect.New())
}
