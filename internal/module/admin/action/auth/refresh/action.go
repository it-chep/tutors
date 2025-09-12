package refresh

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	register_dto "github.com/it-chep/tutors.git/internal/module/admin/action/auth/dto"
	tkn "github.com/it-chep/tutors.git/pkg/token"
)

type Action struct {
	jwtKey, refreshKey string
}

func New(jwtKey, refreshKey string) *Action {
	return &Action{
		jwtKey:     jwtKey,
		refreshKey: refreshKey,
	}
}

func (a *Action) RefreshHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			RefreshToken string `json:"refresh_token"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid request", http.StatusBadRequest)
			return
		}

		token, err := jwt.ParseWithClaims(req.RefreshToken, &register_dto.Claims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(a.refreshKey), nil
		})
		if err != nil || !token.Valid {
			http.Error(w, "invalid refresh token", http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(*register_dto.Claims)
		if !ok || claims.ExpiresAt.Time.Before(time.Now()) {
			http.Error(w, "refresh token expired", http.StatusUnauthorized)
			return
		}

		tokens, err := tkn.GenerateTokens(claims.Email, a.jwtKey, a.refreshKey)
		if err != nil {
			http.Error(w, "server error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(tokens)
	}
}
