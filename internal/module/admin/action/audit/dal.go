package audit

import (
	"context"

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

func (r *Repository) Create(ctx context.Context, entry Entry) error {
	sql := `
		insert into admin_audit (user_id, description, body, action, entity_name, entity_id)
		values ($1, $2, $3::jsonb, $4, $5, $6)
	`

	_, err := r.pool.Exec(
		ctx,
		sql,
		entry.UserID,
		entry.Description,
		entry.Body,
		entry.Action,
		entry.EntityName,
		entry.EntityID,
	)

	return err
}
