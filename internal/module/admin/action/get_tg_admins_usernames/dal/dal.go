package dal

import (
	"context"
	"github.com/it-chep/tutors.git/internal/module/admin/dal/dao"
	"github.com/it-chep/tutors.git/internal/module/admin/dto"

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

// GetAssistantUsernames возвращает доступные тгшки ассистенту
func (r *Repository) GetAssistantUsernames(ctx context.Context, assistantID int64) (dto.TgAdminUsernames, error) {
	sql := `
		select tau.id, tau.name
		from assistant_tgs at2
		    cross join lateral unnest(at2.available_tg_ids) as tg_id
		    join tg_admins_usernames tau on tau.id = tg_id
		where at2.user_id = $1
		  and at2.available_tg_ids is not null
		order by tau.name
	`

	var availableTgs dao.TgAdminUsernames
	err := pgxscan.Select(ctx, r.pool, &availableTgs, sql, assistantID)
	if err != nil {
		return nil, err
	}

	return availableTgs.ToDomain(), nil
}

func (r *Repository) GetUsernames(ctx context.Context, adminID int64) (dto.TgAdminUsernames, error) {
	sql := `
		select id, name from tg_admins_usernames where admin_id = $1 order by name
	`

	var usernames dao.TgAdminUsernames
	err := pgxscan.Select(ctx, r.pool, &usernames, sql, adminID)
	if err != nil {
		return nil, err
	}

	return usernames.ToDomain(), nil
}
