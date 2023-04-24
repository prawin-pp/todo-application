package internal

import "context"

func GetUserIDFromContext(ctx context.Context) string {
	return "MOCK_USER_ID"
}
