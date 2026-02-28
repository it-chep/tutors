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
		SELECT tau.id, tau.name
		FROM assistant_tgs at2
		    CROSS JOIN LATERAL unnest(at2.available_tg_ids) AS tg_id
		    JOIN tg_admins_usernames tau ON tau.id = tg_id
		WHERE at2.user_id = $1
		  AND at2.available_tg_ids IS NOT NULL
		ORDER BY tau.name
	`

	var availableTgs dao.TgAdminUsernames
	err := pgxscan.Select(ctx, r.pool, &availableTgs, sql, assistantID)
	if err != nil {
		return nil, err
	}

	return availableTgs.ToDomain(), nil
}
