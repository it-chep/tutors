package dal

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/pkg/errors"

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
		select full_name from users where tutor_id = $1
	`
	var name string
	err := pgxscan.Get(ctx, r.pool, &name, sql, tutorID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", nil
		}
		return "", err
	}
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

func (r *Repository) GetStudentAdminID(ctx context.Context, studentID int64) (int64, error) {
	sql := `
		select t.admin_id from students s join tutors t on s.tutor_id = t.id where s.id = $1
	`
	var adminID int64
	err := pgxscan.Get(ctx, r.pool, &adminID, sql, studentID)
	if err != nil {
		return 0, err
	}

	return adminID, nil
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
		select count(*) from conducted_lessons where is_trial = false and student_id = $1
	`
	var count int
	err := pgxscan.Get(ctx, r.pool, &count, sql, studentID)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *Repository) GetStudentPayment(ctx context.Context, studentID int64) (dto.Payment, error) {
	sql := `
		select pc.id, pc.bank 
		from payment_cred pc 
		    join students s on pc.id = s.payment_id 
		where s.id = $1
	`

	var payment dao.Payment
	err := pgxscan.Get(ctx, r.pool, &payment, sql, studentID)
	if err != nil {
		return dto.Payment{}, err
	}

	return payment.ToDomain(), nil
}
