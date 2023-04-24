package todotask

import (
	"net/http"

	"github.com/uptrace/bunrouter"
)

func (s *Server) HandleDeleteTask(w http.ResponseWriter, r bunrouter.Request) error {
	w.WriteHeader(http.StatusNoContent)
	return nil
}
