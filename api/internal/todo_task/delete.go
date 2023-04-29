package todotask

import (
	"net/http"

	"github.com/parwin-pp/todo-application/internal"
	"github.com/parwin-pp/todo-application/internal/httperror"
	"github.com/uptrace/bunrouter"
)

func (s *Server) HandleDeleteTask(w http.ResponseWriter, r bunrouter.Request) error {
	userID := internal.UserIDFromContext(r.Context())
	todoID := r.Param("todoId")
	taskID := r.Param("taskId")

	if err := s.db.DeleteTask(r.Context(), userID, todoID, taskID); err != nil {
		return httperror.ErrInternalServer
	}

	w.WriteHeader(http.StatusNoContent)
	return nil
}
