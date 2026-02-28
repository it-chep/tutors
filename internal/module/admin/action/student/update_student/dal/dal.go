package dal

import (
	"context"
	"github.com/it-chep/tutors.git/internal/module/admin/action/student/update_student/dto"
	"github.com/jackc/pgx/v5/pgxpool"
	"strings"
)

type Repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{
		pool: pool,
	}
}

// GetStudentTgInfo получение текущего tg_admin_username_id студента и admin_id
func (r *Repository) GetStudentTgInfo(ctx context.Context, studentID int64) (adminID int64, tgID int64, err error) {
	sql := `
		select t.admin_id, coalesce(s.tg_admin_username_id, 0) 
		from students s 
			join tutors t on s.tutor_id = t.id 
		where s.id = $1
	`
	err = r.pool.QueryRow(ctx, sql, studentID).Scan(&adminID, &tgID)
	return
}

// UpdateStudent обновление студента
func (r *Repository) UpdateStudent(ctx context.Context, studentID int64, update dto.UpdateRequest) error {
	sql := `
		update students set
			first_name = $2,
			last_name = $3,
			middle_name = $4,
			phone = $5,
			tg = $6,
			cost_per_hour = $7,
			parent_full_name = $8,
			parent_phone = $9,
			parent_tg = $10,
			tg_admin_username_id = nullif($11, 0)
		where id = $1
	`
	args := []interface{}{
		studentID,
		strings.TrimSpace(update.FirstName),
		strings.TrimSpace(update.LastName),
		strings.TrimSpace(update.MiddleName),
		update.Phone,
		update.Tg,
		update.CostPerHour,
		strings.TrimSpace(update.ParentFullName),
		update.ParentPhone,
		update.ParentTg,
		update.TgAdminUsernameID,
	}

	_, err := r.pool.Exec(ctx, sql, args...)
	return err
}
