package get_user_info_dal

import (
	"context"
	"encoding/json"
	"github.com/jackc/pgx/v5"
	"github.com/pkg/errors"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/it-chep/tutors.git/internal/module/admin/dal/dao"
	"github.com/it-chep/tutors.git/internal/module/admin/dto"
	"github.com/it-chep/tutors.git/pkg/xo"
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
	sql := "select * from users where id = $1"
	userDao := &dao.User{}
	if err := pgxscan.Get(ctx, r.pool, userDao, sql, userID); err != nil {
		return nil, err
	}

	return userDao.UserInfo(), nil
}

func (r *Repository) GetPaidFunctions(ctx context.Context, adminID int64) (*dto.PaidFunctions, error) {
	sql := "select * from paid_functions where admin_id = $1"
	p := &xo.PaidFunction{}
	if err := pgxscan.Get(ctx, r.pool, p, sql, adminID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	paid := &dto.PaidFunctions{
		AdminID:       p.AdminID,
		PaidFunctions: make(map[string]bool),
	}

	if err := json.Unmarshal(p.Functions, &paid.PaidFunctions); err != nil {
		return nil, err
	}

	return paid, nil
}
