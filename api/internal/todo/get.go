package todo

import (
	"net/http"

	"github.com/parwin-pp/todo-application/internal"
	"github.com/uptrace/bunrouter"
)

func (s *Server) HandleGetTodos(w http.ResponseWriter, r bunrouter.Request) error {
	userID := internal.GetUserIDFromContext(r.Context())

	todos, err := s.db.GetTodos(r.Context(), userID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return nil
	}

	return bunrouter.JSON(w, todos)
}

func (s *Server) HandleGetTodo(w http.ResponseWriter, r bunrouter.Request) error {
	userID := internal.GetUserIDFromContext(r.Context())
	todoID := r.Param("todoId")

	todo, err := s.db.GetTodo(r.Context(), userID, todoID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return nil
	}
	if todo == nil {
		w.WriteHeader(http.StatusNotFound)
		return nil
	}

	return bunrouter.JSON(w, todo)
}
