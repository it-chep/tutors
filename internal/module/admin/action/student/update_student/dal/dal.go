package dal

import (
	"context"

	"github.com/it-chep/tutors.git/internal/module/admin/action/student/update_student/dto"
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

// UpdateStudent создание студента
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
			parent_tg = $10
		where id = $1
	`
	args := []interface{}{
		studentID,
		update.FirstName,
		update.LastName,
		update.MiddleName,
		update.Phone,
		update.Tg,
		update.CostPerHour,
		update.ParentFullName,
		update.ParentPhone,
		update.ParentTg,
	}

	_, err := r.pool.Exec(ctx, sql, args...)
	return err
}
