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

func (r *Repository) GetTutorsByAdmin(ctx context.Context, adminID int64) ([]dto.Tutor, error) {
	sql := `
		select 
            t.cost_per_hour,
            t.subject_id,
            t.admin_id,
            u.full_name as full_name,
            u.tutor_id as id,
            u.tg,
            u.phone
		from tutors t 
		    join users u on t.id = u.tutor_id 
		where t.admin_id = $1
		order by t.id
	`
	var tutors dao.TutorsDao
	if err := pgxscan.Select(ctx, r.pool, &tutors, sql, adminID); err != nil {
		return nil, err
	}

	return tutors.ToDomain(), nil
}

func (r *Repository) GetTutors(ctx context.Context) ([]dto.Tutor, error) {
	sql := `
		select 
            t.cost_per_hour,
            t.subject_id,
            t.admin_id,
            u.full_name as full_name,
            u.tutor_id as id,
            u.tg,
            u.phone
		from tutors t 
		    join users u on t.id = u.tutor_id 
		order by t.id
	`
	var tutors dao.TutorsDao
	if err := pgxscan.Select(ctx, r.pool, &tutors, sql); err != nil {
		return nil, err
	}

	return tutors.ToDomain(), nil
}

func (r *Repository) GetTutorsStudents(ctx context.Context, tutorIDs []int64) ([]dto.StudentWithTransactions, error) {
	sql := `
		select 
            s.id as student_id,
            s.tutor_id,
            s.is_finished_trial,
            COUNT(th.id) as transactions_count,
            w.balance
        from students s 
        join tutors t on s.tutor_id = t.id
        left join transactions_history th on s.id = th.student_id
        left join wallet w on s.id = w.student_id
        where t.id = any($1)
        group by 
            s.id, 
            s.tutor_id, 
            s.is_finished_trial, 
            w.balance
		`

	var students dao.StudentsWithTransactions
	if err := pgxscan.Select(ctx, r.pool, &students, sql, tutorIDs); err != nil {
		return nil, err
	}

	return students.ToDomain(), nil
}
