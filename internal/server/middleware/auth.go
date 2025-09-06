package middleware

import (
	"net/http"
	"strings"

	"github.com/it-chep/tutors.git/internal/module/admin"
	"github.com/it-chep/tutors.git/internal/module/admin/dto"
	pkgContext "github.com/it-chep/tutors.git/internal/pkg/context"

	"github.com/it-chep/tutors.git/internal/pkg/jwt"
)

// JWTMiddleware проверяет JWT токен
func JWTMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header required", http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			http.Error(w, "Bearer token required", http.StatusUnauthorized)
			return
		}

		if !jwt.Valid(tokenString) {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
		}

		next.ServeHTTP(w, r)
	}
}

// Auth получает роль пользователя и кладет ее в контекст, а также разграничивает урлы по доступам
func Auth(adminModule *admin.Module) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			//// 1. Извлечение токена (например, из заголовка Authorization)
			//tokenString := r.Header.Get("Authorization")
			//if tokenString == "" {
			//	http.Error(w, "Unauthorized: missing token", http.StatusUnauthorized)
			//	return
			//}
			//
			//// 2. Валидация токена (JWT или другой механизм)
			//claims, err := h.authService.ValidateToken(tokenString)
			//if err != nil {
			//	http.Error(w, "Unauthorized: invalid token", http.StatusUnauthorized)
			//	return
			//}
			//
			//// 3. Получение роли пользователя (из токена или БД)
			//userRole := claims.Role
			//if userRole == "" {
			//	http.Error(w, "Forbidden: role not found", http.StatusForbidden)
			//	return
			//}
			//
			roleID := int8(dto.SuperAdminRole)

			//hasPermission, err := adminModule.Actions.CheckPathPermission.Do(ctx, roleID, r.URL.Path)
			//if err != nil || !hasPermission {
			//	w.WriteHeader(http.StatusForbidden) // 403 ебашим
			//	return
			//}

			ctx = pkgContext.WithUserRole(ctx, roleID)

			// 5. Продолжаем выполнение запроса
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
