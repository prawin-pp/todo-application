package todotask

import (
	"net/http"

	"github.com/parwin-pp/todo-application/internal"
	"github.com/uptrace/bunrouter"
)

func (s *Server) HandleGetTasks(w http.ResponseWriter, r bunrouter.Request) error {
	uid := internal.GetUserIDFromContext(r.Context())
	todoID := r.Param("todoId")

	tasks, _ := s.db.GetAll(uid, todoID)

	return bunrouter.JSON(w, tasks)
}
