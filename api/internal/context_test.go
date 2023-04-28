package internal

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUserIDFromContext(t *testing.T) {
	t.Run("should return user id = 'MOCK_USER_ID'", func(t *testing.T) {
		ctx := context.Background()
		ctx = context.WithValue(ctx, AuthContextKey{}, "MOCK_USER_ID")

		userID := UserIDFromContext(ctx)

		require.Equal(t, "MOCK_USER_ID", userID)
	})

	t.Run("should return user id = 'MOCK_SOMETHING'", func(t *testing.T) {
		ctx := context.Background()
		ctx = context.WithValue(ctx, AuthContextKey{}, "MOCK_SOMETHING")

		userID := UserIDFromContext(ctx)

		require.Equal(t, "MOCK_SOMETHING", userID)
	})

	t.Run("should return empty user id if not set value in auth context key", func(t *testing.T) {
		ctx := context.Background()

		userID := UserIDFromContext(ctx)

		require.Equal(t, "", userID)
	})
}
