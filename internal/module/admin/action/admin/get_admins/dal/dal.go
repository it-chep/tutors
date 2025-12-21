package dal

import (
	"context"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/it-chep/tutors.git/internal/module/admin/dal/dao"
	"github.com/it-chep/tutors.git/internal/module/admin/dto"

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

// GetAdmins получение админов для суперадминки
func (r *Repository) GetAdmins(ctx context.Context, role dto.Role) ([]dto.User, error) {
	sql := `
		select u.* from users u join roles r on u.role_id = r.id where r.id = $1
	`

	var admins dao.Users
	err := pgxscan.Select(ctx, r.pool, &admins, sql, role)
	if err != nil {
		return nil, err
	}

	return admins.ToDomain(), nil
}
