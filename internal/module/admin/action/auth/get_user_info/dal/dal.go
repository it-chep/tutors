package get_user_info_dal

import (
	"context"

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

func (r *Repository) GetUser(ctx context.Context, userID int64) (*dto.UserInfo, error) {
	userDao := &dao.User{}
	if err := pgxscan.Select(ctx, r.pool, userDao, "select * from users where id = $1", userID); err != nil {
		return nil, err
	}

	return userDao.UserInfo(), nil
}
