package dal

import (
	"context"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/it-chep/tutors.git/internal/module/admin/action/lessons/update_lesson/dto"
	"github.com/it-chep/tutors.git/internal/module/admin/dal/dao"
	indto "github.com/it-chep/tutors.git/internal/module/admin/dto"
	"github.com/it-chep/tutors.git/internal/pkg/transaction/wrapper"
	"github.com/shopspring/decimal"
)

type Repository struct {
	db wrapper.Database
}

func NewRepository(db wrapper.Database) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) GetLessonByID(ctx context.Context, lessonID int64) (indto.Lesson, error) {
	sql := `select * from conducted_lessons where id = $1`

	var lesson dao.LessonDefaultDAO
	err := pgxscan.Get(ctx, r.db.Pool(ctx), &lesson, sql, lessonID)
	if err != nil {
		return indto.Lesson{}, err
	}

	return lesson.ToDomain(), nil
}

func (r *Repository) GetStudentWallet(ctx context.Context, studentID int64) (indto.Wallet, error) {
	sql := `select * from wallet where student_id = $1`

	var wallet dao.Wallet
	err := pgxscan.Get(ctx, r.db.Pool(ctx), &wallet, sql, studentID)
	if err != nil {
		return indto.Wallet{}, err
	}

	return wallet.ToDomain(), nil
}

func (r *Repository) GetStudentInfo(ctx context.Context, studentID int64) (indto.Student, error) {
	sql := `select * from students where id = $1`

	var student dao.StudentDAO
	err := pgxscan.Get(ctx, r.db.Pool(ctx), &student, sql, studentID)
	if err != nil {
		return indto.Student{}, err
	}
	return student.ToDomain(), nil
}

func (r *Repository) GetTutorInfo(ctx context.Context, tutorID int64) (indto.Tutor, error) {
	sql := `
		select
		    t.id,
			t.cost_per_hour,
			t.subject_id,
			t.admin_id,
			u.full_name,
			u.tg,
			u.phone
		from tutors t
		join users u on t.id = u.tutor_id
		where t.id = $1
	`

	var tutor dao.TutorDAO
	err := pgxscan.Get(ctx, r.db.Pool(ctx), &tutor, sql, tutorID)
	if err != nil {
		return indto.Tutor{}, err
	}

	return tutor.ToDomain(), nil
}

func (r *Repository) UpdateLesson(ctx context.Context, lessonID int64, upd dto.UpdateLesson) error {
	sql := `
		update conducted_lessons 
			set duration_in_minutes = $2, created_at = $3 
			where id = $1
	`
	args := []interface{}{
		lessonID,
		upd.Duration.Minutes(),
		upd.Date.UTC(),
	}

	_, err := r.db.Pool(ctx).Exec(ctx, sql, args...)
	return err
}

func (r *Repository) UpdateStudentBalance(ctx context.Context, studentID int64, balance decimal.Decimal) error {
	sql := `
		update wallet 
		set balance = $2 
		where student_id = $1
	`
	args := []interface{}{
		studentID,
		balance,
	}

	_, err := r.db.Pool(ctx).Exec(ctx, sql, args...)
	return err
}
