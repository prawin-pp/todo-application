package todotask

import (
	"context"

	"github.com/parwin-pp/todo-application/internal/model"
)

type Server struct {
	db Database
}

type Database interface {
	GetAll(ctx context.Context, userID, todoID string) ([]model.TodoTask, error)
	Create(ctx context.Context, userID, todoID string, req model.CreateTodoTaskRequest) (*model.TodoTask, error)
	PartialUpdate(ctx context.Context, userID, todoID, taskID string, req model.PartialUpdateTodoTaskRequest) (*model.TodoTask, error)
	Delete(ctx context.Context, userID, todoID, taskID string) error
}

func NewServer(db Database) *Server {
	return &Server{db: db}
}
