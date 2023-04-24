package internal

import (
	"context"

	"github.com/parwin-pp/todo-application/internal/auth"
)

func GetUserIDFromContext(ctx context.Context) string {
	key := ctx.Value(auth.AuthContextKey{})
	uid, _ := key.(string)
	return uid
}
