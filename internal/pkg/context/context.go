package context

import "context"

// contextKey типобезопасный ключ для контекста
type contextKey string

const (
	// UserRoleKey ключ для роли пользователя в контексте
	UserRoleKey contextKey = "user_role"

	// UserIDKey ключ для ID пользователя в контексте
	UserIDKey contextKey = "user_id"
)

// GetUserRole извлекает роль пользователя из контекста
func GetUserRole(ctx context.Context) (int8, bool) {
	role, ok := ctx.Value(UserRoleKey).(int8)
	return role, ok
}

// GetUserID извлекает ID пользователя из контекста
func GetUserID(ctx context.Context) (string, bool) {
	id, ok := ctx.Value(UserIDKey).(string)
	return id, ok
}
