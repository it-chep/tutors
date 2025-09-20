package context

import "context"

type contextKey string

const (
	userIDKey contextKey = "user_id"
)

func CtxWithUser(ctx context.Context, userID int64) context.Context {
	return context.WithValue(ctx, userIDKey, userID)
}

func UserFromContext(ctx context.Context) int64 {
	userID := ctx.Value(userIDKey)
	if userID == nil {
		return 0
	}
	return userID.(int64)
}
