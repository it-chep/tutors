package middleware

import (
	"net/http"

	"github.com/it-chep/tutors.git/internal/pkg/logger"
)

// LoggerMiddleware добавляет логгер в контекст запроса
func LoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		ctx = logger.ContextWithLogger(ctx, logger.New())
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
