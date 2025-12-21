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

func (r *Repository) CreateAdmin(ctx context.Context, createDTO dto.CreateRequest) (int64, error) {
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
		return 0, err
	}

	if createDTO.Role == dto2.AssistantRole {
		_, err = r.pool.Exec(ctx, "update users set admin_id = $1 where id = $1", userCtx.AdminIDFromContext(ctx))
	} else if createDTO.Role == dto2.AdminRole {
		_, err = r.pool.Exec(ctx, "update users set admin_id = $1 where id = $1", id)
	}

	return id, err
}

// AddAvailableTGs добавления тг доступных ассистенту
func (r *Repository) AddAvailableTGs(ctx context.Context, assistantID int64, createDTO dto.CreateRequest) error {
	sql := `
		insert into assistant_tgs (user_id, available_tgs) values ($1, $2) returning id
	`
	args := []interface{}{
		assistantID,
		createDTO.AvailableTGs,
	}

	_, err := r.pool.Exec(ctx, sql, args...)

	return err
}
