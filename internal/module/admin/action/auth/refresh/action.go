package refresh

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/it-chep/tutors.git/internal/config"
	register_dto "github.com/it-chep/tutors.git/internal/module/admin/action/auth/dto"
	tkn "github.com/it-chep/tutors.git/pkg/token"
)

type Action struct {
	jwt config.JwtConfig
}

func New(jwt config.JwtConfig) *Action {
	return &Action{
		jwt: jwt,
	}
}

func (a *Action) RefreshHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie(register_dto.RefreshCookie)
		if err != nil {
			http.Error(w, "invalid refresh token", http.StatusUnauthorized)
			return
		}

		token, err := jwt.ParseWithClaims(cookie.Value, &register_dto.Claims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(a.jwt.RefreshSecret), nil
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

		tokens, err := tkn.GenerateTokens(claims.Email, a.jwt.JwtSecret, a.jwt.RefreshSecret)
		if err != nil {
			http.Error(w, "server error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		http.SetCookie(w, &http.Cookie{
			Name:     register_dto.RefreshCookie,
			Value:    tokens.Refresh(),
			Expires:  time.Now().UTC().Add(14 * 24 * time.Hour),
			HttpOnly: true,
		})
		_ = json.NewEncoder(w).Encode(tokens)
	}
}
