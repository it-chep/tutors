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

func (r *Repository) GetAllStudentsForAdmin(ctx context.Context, adminID int64) ([]dto.Student, error) {
	sql := `
		select s.* from students s join tutors t on s.tutor_id = t.id where t.admin_id = $1
	`
	var students dao.StudentsDAO
	err := pgxscan.Select(ctx, r.pool, &students, sql, adminID)
	if err != nil {
		return nil, err
	}

	return students.ToDomain(), nil
}

func (r *Repository) GetAllStudentsForSuperAdmin(ctx context.Context) ([]dto.Student, error) {
	sql := `
		select * from students
	`
	var students dao.StudentsDAO
	err := pgxscan.Select(ctx, r.pool, &students, sql)
	if err != nil {
		return nil, err
	}

	return students.ToDomain(), nil
}

func (r *Repository) GetTutorStudentsForAdmin(ctx context.Context, adminID, tutorID int64) ([]dto.Student, error) {
	sql := `
		select s.* from students s join tutors t on s.tutor_id = t.id where t.admin_id = $1 and s.tutor_id = $2
	`
	var students dao.StudentsDAO
	err := pgxscan.Select(ctx, r.pool, &students, sql, adminID, tutorID)
	if err != nil {
		return nil, err
	}

	return students.ToDomain(), nil
}

func (r *Repository) GetTutorStudents(ctx context.Context, tutorID int64) ([]dto.Student, error) {
	sql := `
		select * from students where tutor_id = $1
	`
	var students dao.StudentsDAO
	err := pgxscan.Select(ctx, r.pool, &students, sql, tutorID)
	if err != nil {
		return nil, err
	}

	return students.ToDomain(), nil
}

// GetStudentsWalletInfo получение информации о кошельке студента
func (r *Repository) GetStudentsWalletInfo(ctx context.Context, studentIDs []int64) (map[int64]dto.Wallet, error) {
	sql := `
		select * from wallet where student_id = any($1)
	`

	var wallet []dao.Wallet
	err := pgxscan.Select(ctx, r.pool, &wallet, sql, studentIDs)
	if err != nil {
		return nil, err
	}

	studentsWallet := make(map[int64]dto.Wallet, len(wallet))
	for _, w := range wallet {
		studentsWallet[w.StudentID] = w.ToDomain()
	}
	return studentsWallet, nil
}

// HasStudentsPayments у студента есть платные занятия
func (r *Repository) HasStudentsPayments(ctx context.Context, studentIDs []int64) (map[int64]bool, error) {
	sql := `
		select student_id, count(*) > 0 as has_payments
        from transactions_history
        where student_id = any ($1) and confirmed_at is not null
        group by student_id
	`

	type result struct {
		StudentID   int64 `db:"student_id"`
		HasPayments bool  `db:"has_payments"`
	}

	var results []result
	err := pgxscan.Select(ctx, r.pool, &results, sql, studentIDs)
	if err != nil {
		return nil, err
	}

	resultMap := make(map[int64]bool, len(studentIDs))
	for _, row := range results {
		resultMap[row.StudentID] = row.HasPayments
	}

	for _, studentID := range studentIDs {
		if _, exists := resultMap[studentID]; !exists {
			resultMap[studentID] = false
		}
	}
	return resultMap, nil
}
