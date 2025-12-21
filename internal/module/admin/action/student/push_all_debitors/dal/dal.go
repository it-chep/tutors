package dal

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

// GetDebitors .
func (r *Repository) GetDebitors(ctx context.Context) (dto.Students, error) {
	sql := `
		select s.* 
		from students s 
			join wallet w on s.id = w.student_id
		where w.balance < 0 and s.parent_tg_id is not null
	`

	var students dao.StudentsDAO
	err := pgxscan.Select(ctx, r.pool, &students, sql)
	if err != nil {
		return nil, err
	}

	return students.ToDomain(), nil
}

func (r *Repository) AddNotificationsToHistory(ctx context.Context, userID int64, parentTgID int64) error {
	sql := `
		insert into notification_history (user_id, parent_tg_id) values ($1, $2)
	`

	args := []interface{}{
		userID,
		parentTgID,
	}

	_, err := r.pool.Exec(ctx, sql, args...)
	return err
}

func (r *Repository) GetStudentsWallet(ctx context.Context, studentsIDs []int64) ([]dto.Wallet, error) {
	sql := `
		select * from wallet where student_id = any($1)
	`

	var wallets dao.Wallets
	err := pgxscan.Select(ctx, r.pool, &wallets, sql, studentsIDs)
	if err != nil {
		return nil, err
	}

	return wallets.ToDomain(), nil
}
