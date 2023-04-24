package todo

import (
	"net/http"

	"github.com/parwin-pp/todo-application/internal"
	"github.com/uptrace/bunrouter"
)

func (s *Server) HandleGetTodos(w http.ResponseWriter, r bunrouter.Request) error {
	uid := internal.GetUserIDFromContext(r.Context())

	todos, err := s.db.GetAll(uid)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return nil
	}

	return bunrouter.JSON(w, todos)
}
