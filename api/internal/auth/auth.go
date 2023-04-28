package auth

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/uptrace/bunrouter"
)

func (s *Server) HandleLogin(w http.ResponseWriter, r bunrouter.Request) error {
	var body LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return nil
	}

	user, err := s.db.GetUserByUsername(r.Context(), body.Username)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return nil
	}
	if user == nil {
		w.WriteHeader(http.StatusUnauthorized)
		return nil
	}
	if err = s.en.CompareHash(user.Password, body.Password); err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return nil
	}

	jwt, _ := s.en.SignAuthToken(user.ID.String(), map[string]interface{}{})
	expires, err := time.ParseDuration(s.config.ExpireDuration)
	if err != nil {
		log.Fatalf("could not parse duration: %v", err)
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    jwt,
		HttpOnly: true,
		Secure:   true,
		Path:     "/",
		Expires:  time.Now().Add(expires),
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
