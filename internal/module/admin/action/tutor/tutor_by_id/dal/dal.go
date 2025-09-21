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

// GetTutor получение репетитора по ID
func (r *Repository) GetTutor(ctx context.Context, tutorID int64) (dto.Tutor, error) {
	sql := `
		select 
            t.cost_per_hour,
            t.subject_id,
            t.admin_id,
            u.full_name as full_name,
            u.tutor_id as id,
            u.tg,
            u.phone,
		 	s.name as "subject_name"
		from tutors t 
		    join subjects s on t.subject_id = s.id 
			join users u on t.id = u.tutor_id
		where t.id = $1
	`

	args := []interface{}{
		tutorID,
	}

	var tutor dao.TutorWithSubjectName
	err := pgxscan.Get(ctx, r.pool, &tutor, sql, args...)
	if err != nil {
		return dto.Tutor{}, err
	}

	tutorDTO := tutor.ToDomain()
	tutorDTO.SubjectName = tutor.SubjectName.String

	return tutorDTO, nil
}
