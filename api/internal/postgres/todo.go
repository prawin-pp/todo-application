package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/parwin-pp/todo-application/internal/model"
)

func (db *DB) GetTodos(ctx context.Context, userID string) ([]model.Todo, error) {
	todos := []model.Todo{}
	err := db.db.NewSelect().Model(&todos).Where("user_id = ?", userID).Scan(ctx)
	return todos, err
}

func (db *DB) GetTodo(ctx context.Context, userID, todoID string) (*model.Todo, error) {
	var todo model.Todo
	if err := db.db.NewSelect().Model(&todo).
		Where("user_id = ?", userID).
		Where("id = ?", todoID).
		Scan(ctx); err != nil {
		return nil, err
	}
	return &todo, nil
}

func (db *DB) CreateTodo(ctx context.Context, userID, name string) (*model.Todo, error) {
	todo := &model.Todo{
		UserID: uuid.MustParse(userID),
		Name:   name,
	}
	_, err := db.db.NewInsert().Model(todo).Returning("*").Exec(ctx)
	return todo, err
}
