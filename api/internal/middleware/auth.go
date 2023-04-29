package middleware

import (
	"context"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/parwin-pp/todo-application/internal"
	"github.com/uptrace/bunrouter"
)

type Encrypter interface {
	VerifyAuthToken(tokenStr string) (*jwt.Token, *jwt.MapClaims, error)
}

func NewAuthMiddleware(encrypter Encrypter) bunrouter.MiddlewareFunc {
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
			ctx = context.WithValue(ctx, internal.AuthContextKey{}, userID)
			return next(w, r.WithContext(ctx))
		}
	}
}
