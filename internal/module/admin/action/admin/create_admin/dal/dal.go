package dal

import (
	"context"

	indto "github.com/it-chep/tutors.git/internal/module/admin/dto"

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
		insert into users (email, full_name, role_id, phone, tg) values ($1, $2, $3, $4, $5)
	`
	args := []interface{}{
		createDTO.Email,
		createDTO.FullName,
		indto.AdminRole,
		createDTO.Phone,
		createDTO.Tg,
	}
	_, err := r.pool.Exec(ctx, sql, args...)
	return err
}
