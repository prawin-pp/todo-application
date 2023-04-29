package todotask

import (
	"net/http"

	"github.com/parwin-pp/todo-application/internal"
	"github.com/parwin-pp/todo-application/internal/httperror"
	"github.com/uptrace/bunrouter"
)

func (s *Server) HandleGetTasks(w http.ResponseWriter, r bunrouter.Request) error {
	userID := internal.UserIDFromContext(r.Context())
	todoID := r.Param("todoId")

	tasks, err := s.db.GetTasks(r.Context(), userID, todoID)
	if err != nil {
		return httperror.ErrInternalServer
	}

	return bunrouter.JSON(w, tasks)
}
