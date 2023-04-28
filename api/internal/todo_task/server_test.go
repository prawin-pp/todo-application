package todotask

import (
	"context"
	"net/http"

	"github.com/parwin-pp/todo-application/internal"
	"github.com/uptrace/bunrouter"
)

func mockAuthMiddleware(getUserID func() string) bunrouter.MiddlewareFunc {
	return func(next bunrouter.HandlerFunc) bunrouter.HandlerFunc {
		return func(w http.ResponseWriter, req bunrouter.Request) error {
			ctx := req.Context()
			ctx = context.WithValue(ctx, internal.AuthContextKey{}, getUserID())
			return next(w, req.WithContext(ctx))
		}
	}
}
