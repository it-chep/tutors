package dal

import (
	"context"
	userCtx "github.com/it-chep/tutors.git/pkg/context"

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
            t.is_archive,
            t.tg_admin_username,
            u.full_name as full_name,
            u.tutor_id as id,
            u.tg,
            u.phone,
            u.created_at
		from tutors t
		    join users u on t.id = u.tutor_id
		where t.admin_id = $1 and t.is_archive is false
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
            t.is_archive,
            t.tg_admin_username,
            u.full_name as full_name,
            u.tutor_id as id,
            u.tg,
            u.phone,
            u.created_at
		from tutors t
		    join users u on t.id = u.tutor_id
		where t.is_archive is false
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

func (r *Repository) GetTutorsAvailableToAssistance(ctx context.Context, assistantID int64) ([]dto.Tutor, error) {
	sql := `
		SELECT DISTINCT
			t.cost_per_hour,
			t.subject_id,
			t.admin_id,
			t.is_archive,
			t.tg_admin_username,
			u.full_name as full_name,
			u.tutor_id  as id,
			u.tg,
			u.phone,
			u.created_at
		FROM tutors t
			JOIN users u ON t.id = u.tutor_id
		WHERE t.admin_id = $2
		  AND (
			-- Случай 1: У ассистента нет ограничений по TG (пустой массив или нет записи)
			NOT EXISTS (
				SELECT 1 FROM assistant_tgs at
				WHERE at.user_id = $1
				  AND at.available_tgs IS NOT NULL
				  AND array_length(at.available_tgs, 1) > 0
			)
			-- Случай 2: tg_admin_username репетитора входит в список разрешённых TG ассистента
			OR t.tg_admin_username IN (
				SELECT unnest(at.available_tgs)
				FROM assistant_tgs at
				WHERE at.user_id = $1
				  AND at.available_tgs IS NOT NULL
			)
		  )
		 and t.is_archive is false
		ORDER BY t.id
	`
	var tutors dao.TutorsDao
	if err := pgxscan.Select(ctx, r.pool, &tutors, sql, assistantID, userCtx.AdminIDFromContext(ctx)); err != nil {
		return nil, err
	}

	return tutors.ToDomain(), nil
}
