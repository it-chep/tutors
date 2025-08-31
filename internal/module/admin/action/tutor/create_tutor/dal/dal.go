package dal

import (
	"context"
	"github.com/it-chep/tutors.git/internal/module/admin/action/tutor/create_tutor/dto"
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

func (r *Repository) CreateTutor(ctx context.Context, createDTO dto.Request, adminID int64) error {
	sql := `
		insert into tutors (full_name, phone, tg, cost_per_hour, subject_id, admin_id)  
		values ($1, $2, $3, $4, $5, $6)
	`

	args := []interface{}{
		createDTO.FullName,
		createDTO.Phone,
		createDTO.Tg,
		createDTO.CostPerHour,
		createDTO.SubjectID,
		adminID,
	}

	_, err := r.pool.Exec(ctx, sql, args...)
	return err
}
