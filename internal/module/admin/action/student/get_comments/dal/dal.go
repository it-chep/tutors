package dal

import (
	"context"
	"time"

	"github.com/georgysavva/scany/v2/pgxscan"
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

func (r *Repository) GetByStudentID(ctx context.Context, studentID int64) ([]dto.Comment, error) {
	sql := `
		select c.id,
		       c.user_id,
		       c.student_id,
		       c.text,
		       c.created_at,
		       coalesce(u.full_name, '') as author_full_name
		from comments c
		    left join users u on u.id = c.user_id
		where c.student_id = $1
		order by c.created_at desc, c.id desc
	`

	type row struct {
		ID             int64     `db:"id"`
		UserID         int64     `db:"user_id"`
		StudentID      int64     `db:"student_id"`
		Text           string    `db:"text"`
		AuthorFullName string    `db:"author_full_name"`
		CreatedAt      time.Time `db:"created_at"`
	}

	var rows []row
	if err := pgxscan.Select(ctx, r.pool, &rows, sql, studentID); err != nil {
		return nil, err
	}

	comments := make([]dto.Comment, 0, len(rows))
	for _, item := range rows {
		comments = append(comments, dto.Comment{
			ID:             item.ID,
			UserID:         item.UserID,
			StudentID:      item.StudentID,
			Text:           item.Text,
			AuthorFullName: item.AuthorFullName,
			CreatedAt:      item.CreatedAt,
		})
	}

	return comments, nil
}
