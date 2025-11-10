package dal

import (
	"context"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/it-chep/tutors.git/internal/module/admin/action/lessons/update_lesson/dto"
	"github.com/it-chep/tutors.git/internal/module/admin/dal/dao"
	indto "github.com/it-chep/tutors.git/internal/module/admin/dto"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shopspring/decimal"
)

type Repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{
		pool: pool,
	}
}

func (r *Repository) GetLessonByID(ctx context.Context, lessonID int64) (indto.Lesson, error) {
	sql := `select * from conducted_lessons where id = $1`

	var lesson dao.LessonDefaultDAO
	err := pgxscan.Get(ctx, r.pool, &lesson, sql, lessonID)
	if err != nil {
		return indto.Lesson{}, err
	}

	return lesson.ToDomain(), nil
}

func (r *Repository) GetStudentWallet(ctx context.Context, studentID int64) (indto.Wallet, error) {
	sql := `select * from wallet where student_id = $1`

	var wallet dao.Wallet
	err := pgxscan.Get(ctx, r.pool, &wallet, sql, studentID)
	if err != nil {
		return indto.Wallet{}, err
	}

	return wallet.ToDomain(), nil
}

func (r *Repository) GetStudentInfo(ctx context.Context, studentID int64) (indto.Student, error) {
	sql := `select * from students where id = $1`

	var student dao.StudentDAO
	err := pgxscan.Get(ctx, r.pool, &student, sql, studentID)
	if err != nil {
		return indto.Student{}, err
	}
	return student.ToDomain(), nil
}

func (r *Repository) UpdateLesson(ctx context.Context, lessonID int64, upd dto.UpdateLesson) error {
	sql := `
		update conducted_lessons 
			set duration_in_minutes = $2, created_at = $3 
			where id = $1
	`
	args := []interface{}{
		lessonID,
		upd.Duration.Minutes(),
		upd.Date.UTC(),
	}

	_, err := r.pool.Exec(ctx, sql, args...)
	return err
}

func (r *Repository) UpdateStudentBalance(ctx context.Context, studentID int64, balance decimal.Decimal) error {
	sql := `
		update wallet 
		set balance = $2 
		where student_id = $1
	`
	args := []interface{}{
		studentID,
		balance,
	}

	_, err := r.pool.Exec(ctx, sql, args...)
	return err
}
