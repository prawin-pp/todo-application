package todotask

import (
	"context"

	"github.com/parwin-pp/todo-application/internal/model"
)

type Server struct {
	db Database
}

type Database interface {
	GetTasks(ctx context.Context, userID, todoID string) ([]model.TodoTask, error)
	CreateTask(ctx context.Context, userID, todoID string, req model.CreateTodoTaskRequest) (*model.TodoTask, error)
	PartialUpdateTask(ctx context.Context, userID, todoID, taskID string, req model.PartialUpdateTodoTaskRequest) (*model.TodoTask, error)
	DeleteTask(ctx context.Context, userID, todoID, taskID string) error
}

func NewServer(db Database) *Server {
	return &Server{db: db}
}
