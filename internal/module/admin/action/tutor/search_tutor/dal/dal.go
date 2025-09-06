package dal

import (
	"context"
	"fmt"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/it-chep/tutors.git/internal/module/admin/dal/dao"
	"github.com/it-chep/tutors.git/internal/module/admin/dto"
	"github.com/jackc/pgx/v5/pgxpool"
	"strings"
)

type Repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{
		pool: pool,
	}
}

// Search поиск репетитора по фио
func (r *Repository) Search(ctx context.Context, query string) ([]dto.Tutor, error) {
	sql := `
		select * from tutors where full_name ilike $1
	`

	searchQuery := fmt.Sprintf("%%%s%%", strings.TrimSpace(query))
	var tutors dao.TutorsDao
	if err := pgxscan.Select(ctx, r.pool, &tutors, sql, searchQuery); err != nil {
		return nil, err
	}

	return tutors.ToDomain(), nil
}
