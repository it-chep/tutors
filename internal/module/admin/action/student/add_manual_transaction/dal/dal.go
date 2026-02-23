package dal

import (
	"context"
	"time"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/google/uuid"
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

// GetStudentAdminID получение admin_id по студенту
func (r *Repository) GetStudentAdminID(ctx context.Context, studentID int64) (int64, error) {
	sql := `select t.admin_id from students s join tutors t on s.tutor_id = t.id where s.id = $1`
	var adminID int64
	err := pgxscan.Get(ctx, r.pool, &adminID, sql, studentID)
	return adminID, err
}

// AddManualTransaction создание ручной транзакции (сразу подтверждённой)
func (r *Repository) AddManualTransaction(ctx context.Context, studentID int64, amount int64) (uuid.UUID, error) {
	sql := `
		insert into transactions_history (id, student_id, amount, created_at, confirmed_at, is_manual)
		values ($1, $2, $3, $4, $4, true)
		returning id
	`
	id := uuid.New()
	now := time.Now()
	_, err := r.pool.Exec(ctx, sql, id, studentID, amount, now)
	return id, err
}
