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

// GetAdminByID получение админа по id
func (r *Repository) GetAdminByID(ctx context.Context, adminID int64) (dto.User, error) {
	sql := `
		select * from users where id = $1
	`

	var admin dao.User
	err := pgxscan.Get(ctx, r.pool, &admin, sql, adminID)
	if err != nil {
		return dto.User{}, err
	}

	return admin.UserDto(), nil
}
