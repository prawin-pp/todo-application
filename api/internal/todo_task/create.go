package todotask

import (
	"net/http"

	"github.com/uptrace/bunrouter"
)

func (s *Server) HandleCreateTask(w http.ResponseWriter, r bunrouter.Request) error {
	w.WriteHeader(http.StatusCreated)
	return nil
}
