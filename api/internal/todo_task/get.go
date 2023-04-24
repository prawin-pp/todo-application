package todotask

import (
	"net/http"

	"github.com/parwin-pp/todo-application/internal"
	"github.com/uptrace/bunrouter"
)

func (s *Server) HandleGetTasks(w http.ResponseWriter, r bunrouter.Request) error {
	uid := internal.GetUserIDFromContext(r.Context())
	todoID := r.Param("todoId")

	tasks, err := s.db.GetAll(uid, todoID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return nil
	}

	return bunrouter.JSON(w, tasks)
}