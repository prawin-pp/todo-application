package auth

import (
	"net/http"

	"github.com/parwin-pp/todo-application/internal"
	"github.com/parwin-pp/todo-application/internal/httperror"
	"github.com/uptrace/bunrouter"
)

func (s *Server) HandleGetMe(w http.ResponseWriter, r bunrouter.Request) error {
	userID := internal.UserIDFromContext(r.Context())

	user, err := s.db.GetUser(r.Context(), userID)
	if err != nil {
		return httperror.ErrInternalServer
	}
	if user == nil {
		return httperror.ErrUnauthorized
	}

	return bunrouter.JSON(w, user)
}
