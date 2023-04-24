package todotask

import (
	"encoding/json"
	"net/http"

	"github.com/parwin-pp/todo-application/internal"
	"github.com/uptrace/bunrouter"
)

func (s *Server) HandleCreateTask(w http.ResponseWriter, r bunrouter.Request) error {
	userID := internal.GetUserIDFromContext(r.Context())
	todoID := r.Param("todoId")

	var body CreateTodoTaskRequest
	json.NewDecoder(r.Body).Decode(&body)

	s.db.Create(userID, todoID, body)

	w.WriteHeader(http.StatusCreated)
	return nil
}
