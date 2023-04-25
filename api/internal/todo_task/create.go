package todotask

import (
	"encoding/json"
	"net/http"

	"github.com/parwin-pp/todo-application/internal"
	"github.com/parwin-pp/todo-application/internal/model"
	"github.com/uptrace/bunrouter"
)

func (s *Server) HandleCreateTask(w http.ResponseWriter, r bunrouter.Request) error {
	userID := internal.GetUserIDFromContext(r.Context())
	todoID := r.Param("todoId")

	var body model.CreateTodoTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return nil
	}

	task, err := s.db.Create(r.Context(), userID, todoID, body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return nil
	}

	w.WriteHeader(http.StatusCreated)
	return bunrouter.JSON(w, task)
}
