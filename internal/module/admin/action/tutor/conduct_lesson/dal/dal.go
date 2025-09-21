package dal

import (
	"context"

	"github.com/shopspring/decimal"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/it-chep/tutors.git/internal/module/admin/dal/dao"
	"github.com/it-chep/tutors.git/internal/module/admin/dto"
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

func (r *Repository) GetStudent(ctx context.Context, studentID int64) (dto.Student, error) {
	sql := `
		select * from students where id = $1
	`
	var student dao.StudentDAO
	err := pgxscan.Get(ctx, r.pool, &student, sql, studentID)
	if err != nil {
		return dto.Student{}, err
	}

	return student.ToDomain(), nil
}

func (r *Repository) GetTutor(ctx context.Context, tutorID int64) (dto.Tutor, error) {
	sql := `
		select 
		    t.id,
			t.cost_per_hour,
			t.subject_id,
			t.admin_id,
			u.full_name,
			u.tg,
			u.phone 
		from tutors t 
		    join users u on t.id = u.tutor_id 
		where t.id = $1
	`

	args := []interface{}{
		tutorID,
	}

	var tutor dao.TutorDAO
	err := pgxscan.Get(ctx, r.pool, &tutor, sql, args...)
	if err != nil {
		return dto.Tutor{}, err
	}

	return tutor.ToDomain(), nil
}

func (r *Repository) GetStudentWallet(ctx context.Context, studentID int64) (dto.Wallet, error) {
	sql := `
		select * from wallet where student_id = $1
	`
	var wallet dao.Wallet
	err := pgxscan.Get(ctx, r.pool, &wallet, sql, studentID)
	if err != nil {
		return dto.Wallet{}, err
	}
	return wallet.ToDomain(), nil
}

func (r *Repository) UpdateStudentWallet(ctx context.Context, studentID int64, remain decimal.Decimal) error {
	sql := `
		update wallet set balance = $1 where student_id = $2
	`

	args := []interface{}{
		remain,
		studentID,
	}

	_, err := r.pool.Exec(ctx, sql, args...)
	return err
}

// ConductLesson помечаем что урок проведен
func (r *Repository) ConductLesson(ctx context.Context, studentID, tutorID, durationInMinutes int64) error {
	sql := `
		insert into conducted_lessons(student_id, tutor_id, duration_in_minutes, is_trial)
		values ($1, $2, $3, false)
	`

	args := []interface{}{
		studentID,
		tutorID,
		durationInMinutes,
	}

	_, err := r.pool.Exec(ctx, sql, args...)
	return err
}
