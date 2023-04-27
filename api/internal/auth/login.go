package auth

import (
	"encoding/json"
	"net/http"

	"github.com/uptrace/bunrouter"
)

func (s *Server) HandleLogin(w http.ResponseWriter, r bunrouter.Request) error {
	var body LoginRequest
	json.NewDecoder(r.Body).Decode(&body)

	if body.Username != "TEST_USERNAME" || body.Password != "TEST_PASSWORD" {
		w.WriteHeader(http.StatusUnauthorized)
		return nil
	}
	return nil
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `jons:"password"`
}
