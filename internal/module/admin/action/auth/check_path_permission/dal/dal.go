package dal

import (
	"context"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{
		pool: pool,
	}
}

// CheckPathPermission проверяет есть ли у данной роли доступ к данному урлу
func (r *Repository) CheckPathPermission(ctx context.Context, roleID int8, path string) (bool, error) {
	sql := `
		select exists(select 1
					  from roles_permissions rp
							   join permissions p on rp.permission_id = p.id
					  where rp.role_id = $1
						and p.url ilike $2)
	`
	var exists bool
	err := pgxscan.Get(ctx, r.pool, &exists, sql, roleID, path)
	if err != nil {
		return false, err
	}
	return exists, nil
}
