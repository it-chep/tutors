package middleware

import (
	"net/http"
	"net/url"
)

type CorsConfig interface {
	GetAllowedHosts() []string
}

func CORSMiddleware(corsConfig CorsConfig) func(next http.Handler) http.Handler {
	allowedHosts := make(map[string]struct{})
	for _, host := range corsConfig.GetAllowedHosts() {
		allowedHosts[host] = struct{}{}
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")

			if origin != "" {
				if u, err := url.Parse(origin); err == nil {
					host := u.Host

					if _, exists := allowedHosts[host]; exists {
						w.Header().Set("Access-Control-Allow-Origin", origin)
						w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")
						w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
						w.Header().Set("Access-Control-Allow-Credentials", "true")
					}
				}
			}

			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusNoContent)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
