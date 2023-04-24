package todotask

import (
	"encoding/json"
	"net/http"

	"github.com/parwin-pp/todo-application/internal"
	"github.com/uptrace/bunrouter"
)

func (s *Server) HandlePartialUpdateTask(w http.ResponseWriter, r bunrouter.Request) error {
	userID := internal.GetUserIDFromContext(r.Context())
	todoID := r.Param("todoId")
	taskID := r.Param("taskId")

	var body PartialUpdateTodoTaskRequest
	json.NewDecoder(r.Body).Decode(&body)

	s.db.PartialUpdate(userID, todoID, taskID, body)

	return nil
}
