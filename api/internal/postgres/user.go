package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/parwin-pp/todo-application/internal/model"
)

func (db *DB) GetUserByUsername(ctx context.Context, username string) (*model.User, error) {
	var user model.User
	err := db.db.NewSelect().Model(&user).Where("username = ?", username).Scan(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (db *DB) GetUser(ctx context.Context, userID string) (*model.User, error) {
	var user model.User
	err := db.db.NewSelect().Model(&user).Where("id = ?", userID).Scan(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}
