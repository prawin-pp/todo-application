package todo

import (
	"net/http"

	"github.com/uptrace/bunrouter"
)

func (s *Server) HandleGetTodos(w http.ResponseWriter, r bunrouter.Request) error {
	s.db.GetAll("")
	return nil
}
