package dal

import (
	"context"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Dal struct {
	pool *pgxpool.Pool
}

func NewDal(pool *pgxpool.Pool) *Dal {
	return &Dal{
		pool: pool,
	}
}

// IsParentAlreadyAttached TGID уже к кому-то прикреплен ?
func (r *Dal) IsParentAlreadyAttached(ctx context.Context, parentTgID int64) (bool, error) {
	sql := `select exists(select 1 from students where parent_tg_id = $1)`

	var exists bool
	err := pgxscan.Get(ctx, r.pool, &exists, sql, parentTgID)
	if err != nil {
		return false, err
	}
	return exists, nil
}

// IsStudentExist студент есть ?
func (r *Dal) IsStudentExist(ctx context.Context, studentID int64) (bool, error) {
	sql := `select exists(select 1 from students where id = $1)`

	var exists bool
	err := pgxscan.Get(ctx, r.pool, &exists, sql, studentID)
	if err != nil {
		return false, err
	}
	return exists, nil
}

// IsStudentAlreadyWithTgID данный студент уже с ТГ ?
func (r *Dal) IsStudentAlreadyWithTgID(ctx context.Context, studentID int64) (bool, error) {
	sql := `select exists(select 1 from students where id = $1 and parent_tg_id is not null)`

	var exists bool
	err := pgxscan.Get(ctx, r.pool, &exists, sql, studentID)
	if err != nil {
		return false, err
	}
	return exists, nil
}

// AttachParentToStudent прикрепляем ID к студенту
func (r *Dal) AttachParentToStudent(ctx context.Context, studentID, parentTgID int64) error {
	sql := `update students set parent_tg_id = $1 where id = $2`

	_, err := r.pool.Exec(ctx, sql, parentTgID, studentID)
	return err
}
