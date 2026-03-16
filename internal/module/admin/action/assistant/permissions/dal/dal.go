package dal

import (
	"context"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	pool *pgxpool.Pool
}

type Permissions struct {
	CanViewContracts      bool    `db:"can_view_contracts"`
	CanPenalizeAssistants []int64 `db:"can_penalize_assistant_ids"`
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

func (r *Repository) GetAssistantAdminID(ctx context.Context, assistantID int64) (int64, error) {
	sql := `select admin_id from users where id = $1`
	var adminID int64
	err := pgxscan.Get(ctx, r.pool, &adminID, sql, assistantID)
	return adminID, err
}

func (r *Repository) Update(ctx context.Context, assistantID int64, canViewContracts bool, canPenalizeAssistants []int64) error {
	sql := `
		update assistant_tgs
		set can_view_contracts = $2,
			can_penalize_assistant_ids = $3
		where user_id = $1
	`
	_, err := r.pool.Exec(ctx, sql, assistantID, canViewContracts, canPenalizeAssistants)
	return err
}

func (r *Repository) Get(ctx context.Context, assistantID int64) (Permissions, error) {
	sql := `
		select
			can_view_contracts,
			can_penalize_assistant_ids
		from assistant_tgs
		where user_id = $1
	`
	var permissions Permissions
	err := pgxscan.Get(ctx, r.pool, &permissions, sql, assistantID)
	return permissions, err
}
