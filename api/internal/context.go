package internal

import "context"

type AuthContextKey struct{}

func NewContextWithUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, AuthContextKey{}, userID)
}

func UserIDFromContext(ctx context.Context) string {
	key := ctx.Value(AuthContextKey{})
	uid, _ := key.(string)
	return uid
}
