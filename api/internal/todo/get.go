package todo

import (
	"net/http"

	"github.com/parwin-pp/todo-application/internal"
	"github.com/parwin-pp/todo-application/internal/httperror"
	"github.com/uptrace/bunrouter"
)

func (s *Server) HandleGetTodos(w http.ResponseWriter, r bunrouter.Request) error {
	userID := internal.UserIDFromContext(r.Context())

	todos, err := s.db.GetTodos(r.Context(), userID)
	if err != nil {
		return httperror.ErrInternalServer
	}

	return bunrouter.JSON(w, todos)
}

func (s *Server) HandleGetTodo(w http.ResponseWriter, r bunrouter.Request) error {
	userID := internal.UserIDFromContext(r.Context())
	todoID := r.Param("todoId")

	todo, err := s.db.GetTodo(r.Context(), userID, todoID)
	if err != nil {
		return httperror.ErrInternalServer
	}
	if todo == nil {
		return httperror.ErrNotFound.WithMessage("todo not found")
	}

	return bunrouter.JSON(w, todo)
}
