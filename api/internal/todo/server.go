package todo

import (
	"context"

	"github.com/parwin-pp/todo-application/internal/model"
)

type Server struct {
	db Database
}

type Database interface {
	GetAll(ctx context.Context, userID string) ([]model.Todo, error)
	Create(ctx context.Context, userID, name string) (*model.Todo, error)
}

func NewServer(db Database) *Server {
	return &Server{db: db}
}
