package todotask

import (
	"encoding/json"
	"net/http"

	"github.com/parwin-pp/todo-application/internal"
	"github.com/parwin-pp/todo-application/internal/httperror"
	"github.com/parwin-pp/todo-application/internal/model"
	"github.com/uptrace/bunrouter"
)

func (s *Server) HandleCreateTask(w http.ResponseWriter, r bunrouter.Request) error {
	userID := internal.UserIDFromContext(r.Context())
	todoID := r.Param("todoId")

	var body model.CreateTodoTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return httperror.ErrInvalidRequest
	}

	task, err := s.db.CreateTask(r.Context(), userID, todoID, body)
	if err != nil {
		return httperror.ErrInternalServer
	}

	w.WriteHeader(http.StatusCreated)
	return bunrouter.JSON(w, task)
}
