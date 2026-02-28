package dal

import (
	"context"
	"strings"

	"github.com/it-chep/tutors.git/internal/module/admin/action/tutor/update_tutor/dto"
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

// GetTutorTgInfo получение текущего tg_admin_username_id и admin_id репетитора
func (r *Repository) GetTutorTgInfo(ctx context.Context, tutorID int64) (adminID int64, tgID int64, err error) {
	sql := `select admin_id, coalesce(tg_admin_username_id, 0) from tutors where id = $1`
	err = r.pool.QueryRow(ctx, sql, tutorID).Scan(&adminID, &tgID)
	return
}

// UpdateTutor обновление данных репетитора
func (r *Repository) UpdateTutor(ctx context.Context, tutorID int64, upd dto.UpdateRequest) error {
	tutorSQL := `
		update tutors set
			cost_per_hour = $2,
			subject_id = $3,
			tg_admin_username_id = nullif($4, 0)
		where id = $1
	`
	_, err := r.pool.Exec(ctx, tutorSQL, tutorID,
		upd.CostPerHour,
		upd.SubjectID,
		upd.TgAdminUsernameID,
	)
	if err != nil {
		return err
	}

	userSQL := `
		update users set
			full_name = $2,
			phone = $3,
			tg = $4
		where tutor_id = $1
	`
	_, err = r.pool.Exec(ctx, userSQL, tutorID,
		strings.TrimSpace(upd.FullName),
		upd.Phone,
		upd.Tg,
	)
	return err
}
