package todo

import "github.com/parwin-pp/todo-application/internal/model"

type Server struct {
	db Database
}

type Database interface {
	GetAll(userID string) ([]model.Todo, error)
}

func NewServer(db Database) *Server {
	return &Server{db: db}
}
