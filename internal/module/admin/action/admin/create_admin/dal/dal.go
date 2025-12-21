package dal

import (
	"context"

	"github.com/it-chep/tutors.git/internal/module/admin/action/admin/create_admin/dto"
	dto2 "github.com/it-chep/tutors.git/internal/module/admin/dto"
	userCtx "github.com/it-chep/tutors.git/pkg/context"
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
		insert into users (email, full_name, role_id, phone, tg) values ($1, $2, $3, $4, $5) returning id
	`
	args := []interface{}{
		createDTO.Email,
		createDTO.FullName,
		createDTO.Role,
		createDTO.Phone,
		createDTO.Tg,
	}

	var id int64
	err := r.pool.QueryRow(ctx, sql, args...).Scan(&id)
	if err != nil {
		return err
	}

	if createDTO.Role == dto2.AssistantRole {
		_, err = r.pool.Exec(ctx, "update users set admin_id = $1 where id = $1", userCtx.UserIDFromContext(ctx))
	} else if createDTO.Role == dto2.AdminRole {
		_, err = r.pool.Exec(ctx, "update users set admin_id = $1 where id = $1", id)
	}

	return err
}
