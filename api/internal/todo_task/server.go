package todotask

import "github.com/parwin-pp/todo-application/internal/model"

type Server struct {
	db Database
}

type Database interface {
	GetAll(userID, todoID string) ([]model.TodoTask, error)
	Create(userID, todoID string, req CreateTodoTaskRequest) (*model.TodoTask, error)
	Update(userID, todoID string, req UpdateTodoTaskRequest) (*model.TodoTask, error)
	Delete(userID, todoID, taskID string) error
}

func NewServer(db Database) *Server {
	return &Server{db: db}
}
