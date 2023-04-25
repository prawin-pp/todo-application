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
