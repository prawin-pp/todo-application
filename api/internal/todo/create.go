package todo

import (
	"net/http"

	"github.com/uptrace/bunrouter"
)

func (s *Server) HandleCreateTodo(w http.ResponseWriter, r bunrouter.Request) error {
	w.WriteHeader(http.StatusCreated)
	return nil
}
