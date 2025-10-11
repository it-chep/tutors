package dal

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

// DeleteAdmin удаление админа
func (r *Repository) DeleteAdmin(ctx context.Context, adminID int64) error {
	sql := `
		delete from users where id = $1
	`
	_, err := r.pool.Exec(ctx, sql, adminID)
	return err
}

// DeleteTutors удаление репетиторов
func (r *Repository) DeleteTutors(ctx context.Context, adminID int64) error {
	sql := `
		delete from tutors where admin_id = $1
	`

	_, err := r.pool.Exec(ctx, sql, adminID)
	return err
}

// DeleteStudents удаление студентов
func (r *Repository) DeleteStudents(ctx context.Context, adminID int64) error {
	sql := `
		DELETE FROM students 
        USING tutors 
        WHERE students.tutor_id = tutors.id 
        AND tutors.admin_id = $1
	`

	_, err := r.pool.Exec(ctx, sql, adminID)
	return err
}

// DeleteWallets удаление кошельков обновление админов у репетиторов
func (r *Repository) DeleteWallets(ctx context.Context, adminID int64) error {
	sql := `
		DELETE FROM wallet 
		USING students, tutors
		WHERE wallet.student_id = students.id 
		AND students.tutor_id = tutors.id 
		AND tutors.admin_id = $1
	`

	_, err := r.pool.Exec(ctx, sql, adminID)
	return err
}
