package register_dal

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

func (r *Repository) IsEmailExists(ctx context.Context, email string) (bool, error) {
	var exists bool
	err := r.pool.QueryRow(ctx, "select exists (select 1 from users where email=$1 and password='')", email).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}

func (r *Repository) SavePass(ctx context.Context, email, password string) error {
	_, err := r.pool.Exec(ctx, "update users set password=$1 where email=$2 and password = ''", password, email)
	return err
}
