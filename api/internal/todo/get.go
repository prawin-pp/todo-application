package todo

import (
	"net/http"

	"github.com/parwin-pp/todo-application/internal/auth"
	"github.com/uptrace/bunrouter"
)

func (s *Server) HandleGetTodos(w http.ResponseWriter, r bunrouter.Request) error {
	userID := r.Context().Value(auth.AuthContextKey{})
	uid, _ := userID.(string)

	s.db.GetAll(uid)
	return nil
}
