package todo

import (
	"encoding/json"
	"net/http"

	"github.com/parwin-pp/todo-application/internal"
	"github.com/parwin-pp/todo-application/internal/model"
	"github.com/uptrace/bunrouter"
)

func (s *Server) HandleCreateTodo(w http.ResponseWriter, r bunrouter.Request) error {
	var todo model.Todo

	if err := json.NewDecoder(r.Body).Decode(&todo); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return nil
	}

	uid := internal.GetUserIDFromContext(r.Context())

	s.db.Create(uid, todo.Name)

	w.WriteHeader(http.StatusCreated)
	return nil
}
