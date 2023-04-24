package todotask

import (
	"net/http"

	"github.com/parwin-pp/todo-application/internal"
	"github.com/uptrace/bunrouter"
)

func (s *Server) HandleDeleteTask(w http.ResponseWriter, r bunrouter.Request) error {
	userID := internal.GetUserIDFromContext(r.Context())
	todoID := r.Param("todoId")
	taskID := r.Param("taskId")

	if err := s.db.Delete(r.Context(), userID, todoID, taskID); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return nil
	}

	w.WriteHeader(http.StatusNoContent)
	return nil
}
