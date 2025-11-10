package dal

import (
	"context"
	"github.com/georgysavva/scany/v2/pgxscan"
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

func (r *Repository) GetUsernames(ctx context.Context, adminID int64) ([]string, error) {
	sql := `
		select distinct s.tg_admin_username 
		from students s 
		    join tutors t on s.tutor_id = t.id 
		where t.admin_id = $1 and s.tg_admin_username is not null and s.tg_admin_username != ''
	`

	var usernames []string
	err := pgxscan.Select(ctx, r.pool, &usernames, sql, adminID)
	if err != nil {
		return nil, err
	}

	return usernames, nil
}
