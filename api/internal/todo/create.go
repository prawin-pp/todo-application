package todo

import (
	"encoding/json"
	"net/http"

	"github.com/parwin-pp/todo-application/internal"
	"github.com/parwin-pp/todo-application/internal/httperror"
	"github.com/parwin-pp/todo-application/internal/model"
	"github.com/uptrace/bunrouter"
)

func (s *Server) HandleCreateTodo(w http.ResponseWriter, r bunrouter.Request) error {
	userID := internal.UserIDFromContext(r.Context())

	var body model.Todo
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return httperror.ErrInvalidRequest
	}

	todo, err := s.db.CreateTodo(r.Context(), userID, body.Name)
	if err != nil {
		return httperror.ErrInternalServer
	}

	w.WriteHeader(http.StatusCreated)
	return bunrouter.JSON(w, todo)
}
