package check_path_permission

import (
	"context"
	"github.com/it-chep/tutors.git/internal/module/admin/action/auth/check_path_permission/dal"
	"github.com/jackc/pgx/v5/pgxpool"
	"regexp"
	"strings"
)

type Action struct {
	dal *dal.Repository
}

func New(pool *pgxpool.Pool) *Action {
	return &Action{
		dal: dal.NewRepository(pool),
	}
}

// Do проверяет есть ли доступ у роли к урлу
func (a *Action) Do(ctx context.Context, roleID int8, path string) (bool, error) {
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
