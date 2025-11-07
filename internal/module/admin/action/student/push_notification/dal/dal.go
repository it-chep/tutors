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

func (r *Repository) GetStudentByID(ctx context.Context, studentID int64) (dto.Student, error) {
	sql := `select * from students where id = $1`

	var student dao.StudentDAO
	err := pgxscan.Get(ctx, r.pool, &student, sql, studentID)

	return student.ToDomain(), err
}

func (r *Repository) GetStudentWallet(ctx context.Context, studentID int64) (dto.Wallet, error) {
	sql := `
		select * from wallet where student_id = $1
	`
	var wallet dao.Wallet
	err := pgxscan.Get(ctx, r.pool, &wallet, sql, studentID)
	if err != nil {
		return dto.Wallet{}, err
	}
	return wallet.ToDomain(), nil
}

func (r *Repository) AddNotificationToHistory(ctx context.Context, userID, parentTgID int64) error {
	sql := `insert into notification_history (user_id, parent_tg_id) values ($1, $2)`

	args := []interface{}{
		userID,
		parentTgID,
	}

	_, err := r.pool.Exec(ctx, sql, args...)
	return err
}
