package refresh

import (
	"encoding/json"
	"net/http"

	"github.com/it-chep/tutors.git/internal/config"
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
		claims, err := tkn.RefreshClaimsFromRequest(r, a.jwt.RefreshSecret)
		if err != nil {
			http.Error(w, "Кажется, что вам нужно перезайти в лк", http.StatusUnauthorized)
		}

		tokens, err := tkn.GenerateTokens(claims.Email, a.jwt.JwtSecret, a.jwt.RefreshSecret)
		if err != nil {
			http.Error(w, "server error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(tokens)
	}
}
