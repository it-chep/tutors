package middleware

import (
	"net/http"
	"strings"

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
