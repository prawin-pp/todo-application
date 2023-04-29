package auth

import (
	"net/http"

	"github.com/parwin-pp/todo-application/internal"
	"github.com/uptrace/bunrouter"
)

func (s *Server) HandleGetMe(w http.ResponseWriter, r bunrouter.Request) error {
	userID := internal.UserIDFromContext(r.Context())

	user, err := s.db.GetUser(r.Context(), userID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return nil
	}
	if user == nil {
		w.WriteHeader(http.StatusUnauthorized)
		return nil
	}

	return bunrouter.JSON(w, user)
}
