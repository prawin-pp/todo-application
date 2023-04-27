package middleware

import (
	"context"
	"net/http"

	"github.com/parwin-pp/todo-application/internal/auth"
	"github.com/uptrace/bunrouter"
)

func NewAuthMiddleware(encrypter *auth.AuthEncryption) bunrouter.MiddlewareFunc {
	return func(next bunrouter.HandlerFunc) bunrouter.HandlerFunc {
		return func(w http.ResponseWriter, r bunrouter.Request) error {
			token, err := r.Cookie("token")
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				return nil
			}

			_, claims, err := encrypter.VerifyAuthToken(token.Value)
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				return nil
			}

			userID, _ := claims.GetSubject()
			ctx := r.Context()
			ctx = context.WithValue(ctx, auth.AuthContextKey{}, userID)
			return next(w, r.WithContext(ctx))
		}
	}
}
