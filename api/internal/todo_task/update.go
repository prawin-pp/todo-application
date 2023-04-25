package todotask

import (
	"encoding/json"
	"net/http"

	"github.com/parwin-pp/todo-application/internal"
	"github.com/parwin-pp/todo-application/internal/model"
	"github.com/uptrace/bunrouter"
)

func (s *Server) HandlePartialUpdateTask(w http.ResponseWriter, r bunrouter.Request) error {
	userID := internal.GetUserIDFromContext(r.Context())
	todoID := r.Param("todoId")
	taskID := r.Param("taskId")

	var body model.PartialUpdateTodoTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return nil
	}

	updatedTask, err := s.db.PartialUpdateTask(r.Context(), userID, todoID, taskID, body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return nil
	}

	return bunrouter.JSON(w, updatedTask)
}
