package login_dal

import (
	"context"

	"github.com/georgysavva/scany/v2/pgxscan"
	register_dto "github.com/it-chep/tutors.git/internal/module/admin/action/auth/dto"
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

func (r *Repository) GetUser(ctx context.Context, email string) (user *register_dto.User, _ error) {
	sql := `select * from users where email = $1`
	user = &register_dto.User{}

	return user, pgxscan.Get(ctx, r.pool, user, sql, email)
}

func (r *Repository) SaveCode(ctx context.Context, email, code string) error {
	_, err := r.pool.Exec(ctx, "update users set smtp_code=$1 where email=$2", code, email)
	return err
}
