package dal

import (
	"context"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/it-chep/tutors.git/internal/module/admin/dal/dao"
	"github.com/it-chep/tutors.git/internal/module/admin/dto"
	userCtx "github.com/it-chep/tutors.git/pkg/context"
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

func (r *Repository) GetAllTutorsForAdmin(ctx context.Context, adminID int64) ([]dto.Tutor, error) {
	sql := `
		select
			t.cost_per_hour,
			t.subject_id,
			t.admin_id,
			t.is_archive,
			t.tg_admin_username_id,
			tau.name as tg_admin_username,
			u.full_name as full_name,
			u.tutor_id as id,
			u.tg,
			u.phone,
			u.created_at
		from tutors t
		    join users u on t.id = u.tutor_id
		    left join tg_admins_usernames tau on t.tg_admin_username_id = tau.id
		where t.admin_id = $1 and t.is_archive is true
		order by t.id
	`
	var tutors dao.TutorsDao
	if err := pgxscan.Select(ctx, r.pool, &tutors, sql, adminID); err != nil {
		return nil, err
	}

	return tutors.ToDomain(), nil
}

func (r *Repository) GetAllTutorsForSuperAdmin(ctx context.Context) ([]dto.Tutor, error) {
	sql := `
		select
			t.cost_per_hour,
			t.subject_id,
			t.admin_id,
			t.is_archive,
			t.tg_admin_username_id,
			tau.name as tg_admin_username,
			u.full_name as full_name,
			u.tutor_id as id,
			u.tg,
			u.phone,
			u.created_at
		from tutors t
		    join users u on t.id = u.tutor_id
		    left join tg_admins_usernames tau on t.tg_admin_username_id = tau.id
		where t.is_archive is true
		order by t.id
	`
	var tutors dao.TutorsDao
	if err := pgxscan.Select(ctx, r.pool, &tutors, sql); err != nil {
		return nil, err
	}

	return tutors.ToDomain(), nil
}

func (r *Repository) GetTutorsAvailableToAssistant(ctx context.Context, assistantID int64) ([]dto.Tutor, error) {
	sql := `
		select distinct 
			t.cost_per_hour,
			t.subject_id,
			t.admin_id,
			t.is_archive,
			t.tg_admin_username_id,
			u.full_name as full_name,
			u.tutor_id  as id,
			u.tg,
			u.phone,
			u.created_at
		from tutors t
			join users u ON t.id = u.tutor_id
		WHERE t.admin_id = $2
		  AND t.is_archive is true
		  AND (
			NOT EXISTS (
				SELECT 1 FROM assistant_tgs at
				WHERE at.user_id = $1
				  AND at.available_tg_ids IS NOT NULL
				  AND array_length(at.available_tg_ids, 1) > 0
			)
			OR t.tg_admin_username_id IN (
				SELECT unnest(at.available_tg_ids)
				FROM assistant_tgs at
				WHERE at.user_id = $1
				  AND at.available_tg_ids IS NOT NULL
			)
		  )
		ORDER BY u.tutor_id
	`
	var tutors dao.TutorsDao
	if err := pgxscan.Select(ctx, r.pool, &tutors, sql, assistantID, userCtx.AdminIDFromContext(ctx)); err != nil {
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
