package dal

import (
	"context"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/it-chep/tutors.git/internal/module/admin/dal/dao"
	"github.com/it-chep/tutors.git/internal/module/admin/dto"
	userCtx "github.com/it-chep/tutors.git/pkg/context"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

type Repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{
		pool: pool,
	}
}

func (r *Repository) GetAllLessons(ctx context.Context, adminID int64, from, to time.Time) ([]dto.Lesson, error) {
	sql := `
		select cl.*, s.first_name, s.last_name, s.middle_name, u.full_name as tutor_name
		from conducted_lessons cl 
    		join tutors t on cl.tutor_id = t.id 
			join students s on cl.student_id = s.id
    		join users u on t.id = u.tutor_id
		where cl.is_trial is not true and t.admin_id = $1 and cl.created_at between $2 and $3
	`

	args := []interface{}{
		adminID,
		from,
		to,
	}

	if dto.IsAssistantRole(ctx) {
		sql = `
		select cl.*, s.first_name, s.last_name, s.middle_name, u.full_name as tutor_name
		from conducted_lessons cl 
    		join tutors t on cl.tutor_id = t.id 
			join students s on cl.student_id = s.id
    		join users u on t.id = u.tutor_id
		where cl.is_trial is not true and t.admin_id = $1 and cl.created_at between $2 and $3
				and (
                    not exists (
                        select 1 
                        from assistant_tgs at
                        where at.user_id = $4
                          and at.available_tgs is not null
                          and array_length(at.available_tgs, 1) > 0
                    )
                    or s.tg_admin_username in (
                        select unnest(at.available_tgs)
                        from assistant_tgs at
                        where at.user_id = $4
                          and at.available_tgs is not null
                    )
                )
		`
		args = append(args, userCtx.UserIDFromContext(ctx))
	}

	var lessons dao.LessonsDAO
	err := pgxscan.Select(ctx, r.pool, &lessons, sql, args...)
	if err != nil {
		return nil, err
	}

	return lessons.ToDomain(ctx), nil
}
