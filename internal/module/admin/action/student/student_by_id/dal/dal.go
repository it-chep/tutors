package dal

import (
	"context"

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

func (r *Repository) GetSubjectName(ctx context.Context, subjectID int64) (string, error) {
	sql := `
		select name from subjects where id = $1
	`
	var name string
	err := pgxscan.Get(ctx, r.pool, &name, sql, subjectID)
	return name, err
}

func (r *Repository) GetTutorName(ctx context.Context, tutorID int64) (string, error) {
	sql := `
		select full_name from tutors where id = $1
	`
	var name string
	err := pgxscan.Get(ctx, r.pool, &name, sql, tutorID)
	return name, err
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

// GetStudentWalletInfo получение информации о кошельке студента
func (r *Repository) GetStudentWalletInfo(ctx context.Context, studentID int64) (dto.Wallet, error) {
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

// HasStudentPayments у студента есть платные занятия
func (r *Repository) HasStudentPayments(ctx context.Context, studentID int64) (bool, error) {
	sql := `
		select count(*) from transactions_history where student_id = $1
	`
	var count int
	err := pgxscan.Get(ctx, r.pool, &count, sql, studentID)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
