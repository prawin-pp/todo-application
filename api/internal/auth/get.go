package auth

import (
	"net/http"

	"github.com/uptrace/bunrouter"
)

func (s *Server) HandleGetMe(w http.ResponseWriter, r bunrouter.Request) error {
	value := r.Context().Value(AuthContextKey{})

	userID, ok := value.(string)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return nil
	}

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
