package mock

import (
	"context"

	"github.com/parwin-pp/todo-application/internal/model"
)

type TodoDatabase struct {
	GetTodosFn   func(ctx context.Context, userID string) ([]model.Todo, error)
	GetTodoFn    func(ctx context.Context, userID, todoID string) (*model.Todo, error)
	CreateTodoFn func(ctx context.Context, userID, name string) (*model.Todo, error)
}

func (db *TodoDatabase) GetTodos(ctx context.Context, userID string) ([]model.Todo, error) {
	return db.GetTodosFn(ctx, userID)
}

func (db *TodoDatabase) GetTodo(ctx context.Context, userID, todoID string) (*model.Todo, error) {
	return db.GetTodoFn(ctx, userID, todoID)
}

func (db *TodoDatabase) CreateTodo(ctx context.Context, userID, name string) (*model.Todo, error) {
	return db.CreateTodoFn(ctx, userID, name)
}
