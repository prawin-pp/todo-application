package todotask

import (
	"encoding/json"
	"net/http"

	"github.com/parwin-pp/todo-application/internal"
	"github.com/parwin-pp/todo-application/internal/httperror"
	"github.com/parwin-pp/todo-application/internal/model"
	"github.com/uptrace/bunrouter"
)

func (s *Server) HandlePartialUpdateTask(w http.ResponseWriter, r bunrouter.Request) error {
	userID := internal.UserIDFromContext(r.Context())
	todoID := r.Param("todoId")
	taskID := r.Param("taskId")

	var body model.PartialUpdateTodoTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return httperror.ErrInvalidRequest
	}

	updatedTask, err := s.db.PartialUpdateTask(r.Context(), userID, todoID, taskID, body)
	if err != nil {
		return httperror.ErrInternalServer
	}

	return bunrouter.JSON(w, updatedTask)
}
