package dal

import (
	"context"
	"github.com/it-chep/tutors.git/internal/module/admin/action/admin/create_admin/dto"
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

func (r *Repository) CreateAdmin(ctx context.Context, createDTO dto.CreateRequest) error {
	sql := `
		insert into admins (full_name, phone, tg) values ($1, $2, $3)
	`
	args := []interface{}{
		createDTO.FullName,
		createDTO.Phone,
		createDTO.Tg,
	}
	_, err := r.pool.Exec(ctx, sql, args...)
	return err
}
