package dal

import (
	"context"
	"fmt"
	"strings"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/it-chep/tutors.git/internal/module/admin/action/tutor/archive_filter/dto"
	"github.com/it-chep/tutors.git/internal/module/admin/dal/dao"
	indto "github.com/it-chep/tutors.git/internal/module/admin/dto"
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

func (r *Repository) FilterTutors(ctx context.Context, adminID int64, filter dto.FilterRequest) ([]indto.Tutor, error) {
	sql, phValues := stmtBuilder(ctx, adminID, filter)

	var tutors dao.TutorsDao
	err := pgxscan.Select(ctx, r.pool, &tutors, sql, phValues...)
	if err != nil {
		return nil, err
	}

	return tutors.ToDomain(), nil
}

func stmtBuilder(ctx context.Context, adminID int64, filter dto.FilterRequest) (_ string, phValues []any) {
	defaultSql := `
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
	`

	phValues = append(phValues, adminID)

	whereStmtBuilder := strings.Builder{}
	phCounter := 2

	if len(filter.TgUsernameIDs) != 0 {
		whereStmtBuilder.WriteString(
			fmt.Sprintf(`
				and t.tg_admin_username_id = any($%d)
			`, phCounter),
		)
		phValues = append(phValues, filter.TgUsernameIDs)
		phCounter++
	}

	if indto.IsAssistantRole(ctx) {
		assistantID := userCtx.UserIDFromContext(ctx)

		whereStmtBuilder.WriteString(
			fmt.Sprintf(`
                and (
                    not exists (
                        select 1
                        from assistant_tgs at
                        where at.user_id = $%d
                          and at.available_tg_ids is not null
                          and array_length(at.available_tg_ids, 1) > 0
                    )
                    or t.tg_admin_username_id in (
                        select unnest(at.available_tg_ids)
                        from assistant_tgs at
                        where at.user_id = $%d
                          and at.available_tg_ids is not null
                    )
                )
            `, phCounter, phCounter),
		)

		phValues = append(phValues, assistantID)
		phCounter++
	}

	return fmt.Sprintf(`
		%s
		%s
        order by t.id
    `, defaultSql, whereStmtBuilder.String()), phValues
}

func (r *Repository) GetTutorsStudents(ctx context.Context, tutorIDs []int64) ([]indto.StudentWithTransactions, error) {
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
