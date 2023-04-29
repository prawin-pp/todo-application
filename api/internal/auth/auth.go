package auth

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/parwin-pp/todo-application/internal/httperror"
	"github.com/uptrace/bunrouter"
)

func (s *Server) HandleLogin(w http.ResponseWriter, r bunrouter.Request) error {
	var body LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return httperror.ErrInvalidRequest
	}

	user, err := s.db.GetUserByUsername(r.Context(), body.Username)
	if err != nil {
		return httperror.ErrInternalServer
	}
	if user == nil {
		return httperror.ErrUnauthorized
	}
	if err = s.en.CompareHash(user.Password, body.Password); err != nil {
		return httperror.ErrUnauthorized
	}

	jwt, err := s.en.SignAuthToken(user.ID.String(), map[string]interface{}{})
	if err != nil {
		return httperror.ErrInternalServer
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    jwt,
		HttpOnly: true,
		Secure:   true,
		Path:     "/",
		Expires:  time.Now().Add(s.config.ExpireDuration),
	})
	return bunrouter.JSON(w, user)
}

func (s *Server) HandleLogout(w http.ResponseWriter, r bunrouter.Request) error {
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    "",
		HttpOnly: true,
		Secure:   true,
		Path:     "/",
		Expires:  time.Now().Add(-time.Hour),
	})
	return nil
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `jons:"password"`
}
