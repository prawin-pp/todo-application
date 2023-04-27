package auth

import (
	"encoding/json"
	"net/http"

	"github.com/uptrace/bunrouter"
)

func (s *Server) HandleLogin(w http.ResponseWriter, r bunrouter.Request) error {
	var body LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return nil
	}

	if body.Username != "TEST_USERNAME" || body.Password != "TEST_PASSWORD" {
		w.WriteHeader(http.StatusUnauthorized)
		return nil
	}

	s.db.GetUserByUsername(r.Context(), body.Username)
	return nil
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `jons:"password"`
}
