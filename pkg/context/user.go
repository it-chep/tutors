package context

import "context"

type contextKey string

const (
	// userIDKey ключ ID пользователя
	userIDKey contextKey = "user_id"
	// userRoleKey ключ для роли пользователя в контексте
	userRoleKey contextKey = "user_role"
	// tutorIDKey ключ для ID репетитора
	tutorIDKey contextKey = "tutor_id"
)

func WithUserID(ctx context.Context, userID int64) context.Context {
	return context.WithValue(ctx, userIDKey, userID)
}

func UserIDFromContext(ctx context.Context) int64 {
	userID := ctx.Value(userIDKey)
	if userID == nil {
		return 0
	}
	return userID.(int64)
}

// GetUserRole извлекает роль пользователя из контекста
func GetUserRole(ctx context.Context) (int8, bool) {
	role, ok := ctx.Value(userRoleKey).(int8)
	return role, ok
}

// WithUserRole добавляет информацию по роли в контекст
func WithUserRole(ctx context.Context, roleID int8) context.Context {
	return context.WithValue(ctx, userRoleKey, roleID)
}

// WithTutorID добавляет информацию по ID репетитора
func WithTutorID(ctx context.Context, tutorID int64) context.Context {
	return context.WithValue(ctx, tutorIDKey, tutorID)
}

// GetTutorID извлекает id репетитора из контекста
func GetTutorID(ctx context.Context) int64 {
	userID := ctx.Value(tutorIDKey)
	if userID == nil {
		return 0
	}
	return userID.(int64)
}
