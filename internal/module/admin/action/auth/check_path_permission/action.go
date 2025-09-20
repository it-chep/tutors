package check_path_permission

import (
	"context"
	"net/http"
	"regexp"
	"strings"

	"github.com/it-chep/tutors.git/internal/config"
	"github.com/it-chep/tutors.git/internal/module/admin/action/auth/check_path_permission/dal"
	userCtx "github.com/it-chep/tutors.git/pkg/context"
	"github.com/it-chep/tutors.git/pkg/token"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Action struct {
	dal *dal.Repository
	jwt config.JwtConfig
}

func New(pool *pgxpool.Pool, jwt config.JwtConfig) *Action {
	return &Action{
		dal: dal.NewRepository(pool),
		jwt: jwt,
	}
}

func (a *Action) AuthMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			claims, err := token.AccessClaimsFromRequest(r, a.jwt.JwtSecret)
			if err != nil {
				http.Error(w, "authorization required", http.StatusUnauthorized)
			}

			user, err := a.dal.GetUser(ctx, claims.Email)
			if err != nil {
				http.Error(w, "authorization required", http.StatusUnauthorized)
			}

			has, err := a.hasPermissions(ctx, int8(user.Role), r.URL.Path)
			if err != nil || !has {
				http.Error(w, "authorization required", http.StatusUnauthorized)
			}

			ctx = userCtx.CtxWithUser(ctx, user.ID)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// Do проверяет есть ли доступ у роли к урлу
func (a *Action) hasPermissions(ctx context.Context, roleID int8, path string) (bool, error) {
	re := regexp.MustCompile(`/\d+`)

	// Заменяем все числа в пути
	normalized := re.ReplaceAllString(path, `/{id}`)

	// Если URL заканчивается на число (например, "/tutors/1")
	if strings.HasSuffix(normalized, "/{id") && !strings.HasSuffix(normalized, "/{id}") {
		normalized += "}"
	}

	hasPermission, err := a.dal.CheckPathPermission(ctx, roleID, normalized)
	if err != nil {
		return false, err
	}

	return hasPermission, nil
}
