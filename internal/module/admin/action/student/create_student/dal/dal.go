package dal

import (
	"context"

	"github.com/it-chep/tutors.git/internal/module/admin/action/student/create_student/dto"
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

// CreateStudent создание студента
func (r *Repository) CreateStudent(ctx context.Context, createDTO dto.CreateRequest) (int64, error) {
	sql := `
		insert into students (first_name, last_name, middle_name, phone, tg, cost_per_hour, subject_id, tutor_id, is_finished_trial, parent_full_name, parent_phone, parent_tg, tg_admin_username) 
		values ($1, $2, $3, $4, $5, $6, $7, $8, false, $9, $10, $11, $12)
		returning id
	`
	args := []interface{}{
		createDTO.FirstName,
		createDTO.LastName,
		createDTO.MiddleName,
		createDTO.Phone,
		createDTO.Tg,
		createDTO.CostPerHour,
		createDTO.SubjectID,
		createDTO.TutorID,
		createDTO.ParentFullName,
		createDTO.ParentPhone,
		createDTO.ParentTg,
		createDTO.TgAdminUsername,
	}

	var id int64
	row := r.pool.QueryRow(ctx, sql, args...)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

// CreateWallet создание кошелька
func (r *Repository) CreateWallet(ctx context.Context, studentID int64) error {
	sql := `
		insert into wallet ( student_id, balance ) values ( $1, 0 ) 
	`
	args := []interface{}{
		studentID,
	}
	_, err := r.pool.Exec(ctx, sql, args...)
	return err
}

// AddTgToAssistant добавление тг ассистенту
func (r *Repository) AddTgToAssistant(ctx context.Context, assistantID int64, tgAdminUsername string) error {
	sql := `
		update assistant_tgs
		set available_tgs = array(
			select distinct unnest(array_append(available_tgs, $2))
		)
		where user_id = $1
	`
	args := []interface{}{
		assistantID,
		tgAdminUsername,
	}

	_, err := r.pool.Exec(ctx, sql, args...)
	if err != nil {
		return err
	}
	return nil
}
