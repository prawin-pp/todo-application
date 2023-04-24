package todo

import (
	"encoding/json"
	"net/http"

	"github.com/parwin-pp/todo-application/internal"
	"github.com/uptrace/bunrouter"
)

func (s *Server) HandleGetTodos(w http.ResponseWriter, r bunrouter.Request) error {
	uid := internal.GetUserIDFromContext(r.Context())

	todos, _ := s.db.GetAll(uid)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	return json.NewEncoder(w).Encode(todos)
}
