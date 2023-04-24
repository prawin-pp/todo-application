package internal

import (
	"context"
	"testing"

	"github.com/parwin-pp/todo-application/internal/auth"
	"github.com/stretchr/testify/require"
)

func TestGetUserIDFromContext(t *testing.T) {
	t.Run("should return user id = 'MOCK_USER_ID'", func(t *testing.T) {
		ctx := context.Background()
		ctx = context.WithValue(ctx, auth.AuthContextKey{}, "MOCK_USER_ID")

		userID := GetUserIDFromContext(ctx)

		require.Equal(t, "MOCK_USER_ID", userID)
	})

	t.Run("should return user id = 'MOCK_SOMETHING'", func(t *testing.T) {
		ctx := context.Background()
		ctx = context.WithValue(ctx, auth.AuthContextKey{}, "MOCK_SOMETHING")

		userID := GetUserIDFromContext(ctx)

		require.Equal(t, "MOCK_SOMETHING", userID)
	})

}
