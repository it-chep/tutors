package dal

import (
	"context"

	indto "github.com/it-chep/tutors.git/internal/module/admin/dto"

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

// CreateTutor создание репетитора
func (r *Repository) CreateTutor(ctx context.Context, createDTO dto.Request, adminID int64) (int64, error) {
	sql := `
		insert into tutors (phone, tg, cost_per_hour, subject_id, admin_id)  
		values ($1, $2, $3, $4, $5)
		returning id;
	`

	args := []interface{}{
		createDTO.Phone,
		createDTO.Tg,
		createDTO.CostPerHour,
		createDTO.SubjectID,
		adminID,
	}

	var id int64
	err := r.pool.QueryRow(ctx, sql, args...).Scan(&id)
	return id, err
}

// CreateUser создание пользователя
func (r *Repository) CreateUser(ctx context.Context, createDTO dto.Request, tutorID int64) error {
	sql := `
		insert into users (email, password, full_name, is_active, role_id, tutor_id)
		values ($1, $2, $3, false, $4)
	`
	args := []interface{}{
		createDTO.Email,
		// password ?
		createDTO.FullName,
		indto.TutorRole,
		tutorID,
	}
	_, err := r.pool.Exec(ctx, sql, args...)
	return err
}
