package todo

import (
	"encoding/json"
	"net/http"

	"github.com/parwin-pp/todo-application/internal"
	"github.com/parwin-pp/todo-application/internal/model"
	"github.com/uptrace/bunrouter"
)

func (s *Server) HandleCreateTodo(w http.ResponseWriter, r bunrouter.Request) error {
	var body model.Todo

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return nil
	}

	uid := internal.GetUserIDFromContext(r.Context())

	todo, err := s.db.Create(uid, body.Name)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return nil
	}

	w.WriteHeader(http.StatusCreated)
	return bunrouter.JSON(w, todo)
}
