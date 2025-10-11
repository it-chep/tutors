package action

import (
	"github.com/it-chep/tutors.git/internal/config"
	"github.com/it-chep/tutors.git/internal/module/admin/action/admin/create_admin"
	"github.com/it-chep/tutors.git/internal/module/admin/action/admin/delete_admin"
	"github.com/it-chep/tutors.git/internal/module/admin/action/admin/get_admin_by_id"
	"github.com/it-chep/tutors.git/internal/module/admin/action/admin/get_admins"
	"github.com/it-chep/tutors.git/internal/module/admin/action/auth"
	"github.com/it-chep/tutors.git/internal/module/admin/action/get_all_finance"
	"github.com/it-chep/tutors.git/internal/module/admin/action/get_all_subjects"
	"github.com/it-chep/tutors.git/internal/module/admin/action/student/create_student"
	"github.com/it-chep/tutors.git/internal/module/admin/action/student/delete_student"
	"github.com/it-chep/tutors.git/internal/module/admin/action/student/get_student_finance"
	"github.com/it-chep/tutors.git/internal/module/admin/action/student/get_students"
	"github.com/it-chep/tutors.git/internal/module/admin/action/student/move_students"
	"github.com/it-chep/tutors.git/internal/module/admin/action/student/search_student"
	"github.com/it-chep/tutors.git/internal/module/admin/action/student/student_by_id"
	"github.com/it-chep/tutors.git/internal/module/admin/action/tutor/conduct_lesson"
	"github.com/it-chep/tutors.git/internal/module/admin/action/tutor/conduct_trial"
	"github.com/it-chep/tutors.git/internal/module/admin/action/tutor/create_tutor"
	"github.com/it-chep/tutors.git/internal/module/admin/action/tutor/delete_tutor"
	"github.com/it-chep/tutors.git/internal/module/admin/action/tutor/get_tutor_finance"
	"github.com/it-chep/tutors.git/internal/module/admin/action/tutor/get_tutors"
	"github.com/it-chep/tutors.git/internal/module/admin/action/tutor/search_tutor"
	"github.com/it-chep/tutors.git/internal/module/admin/action/tutor/tutor_by_id"
	"github.com/it-chep/tutors.git/internal/pkg/tg_bot"
	"github.com/it-chep/tutors.git/pkg/smtp"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Aggregator struct {
	// Репетитор
	CreateTutor     *create_tutor.Action
	DeleteTutor     *delete_tutor.Action
	GetTutorFinance *get_tutor_finance.Action
	GetTutors       *get_tutors.Action
	SearchTutor     *search_tutor.Action
	GetTutorByID    *tutor_by_id.Action
	ConductTrial    *conduct_trial.Action
	ConductLesson   *conduct_lesson.Action

	// Студент
	CreateStudent     *create_student.Action
	DeleteStudent     *delete_student.Action
	GetStudentFinance *get_student_finance.Action
	GetStudents       *get_students.Action
	SearchStudent     *search_student.Action
	GetStudentByID    *student_by_id.Action
	MoveStudents      *move_students.Action

	// Финансы
	GetAllFinance *get_all_finance.Action

	// Предметы
	GetAllSubjects *get_all_subjects.Action

	// AUTH
	Auth *auth.Aggregator

	// Админы
	CreateAdmin  *create_admin.Action
	DeleteAdmin  *delete_admin.Action
	GetAdmins    *get_admins.Action
	GetAdminByID *get_admin_by_id.Action
}

func NewAggregator(pool *pgxpool.Pool, smtp *smtp.ClientSmtp, config config.JwtConfig, bot *tg_bot.Bot) *Aggregator {
	return &Aggregator{
		// Репетитор
		CreateTutor:     create_tutor.New(pool),
		DeleteTutor:     delete_tutor.New(pool),
		GetTutorFinance: get_tutor_finance.New(pool),
		GetTutors:       get_tutors.New(pool),
		SearchTutor:     search_tutor.New(pool),
		GetTutorByID:    tutor_by_id.New(pool),
		ConductTrial:    conduct_trial.New(pool),
		ConductLesson:   conduct_lesson.New(pool, bot),

		// Студент
		CreateStudent:     create_student.New(pool),
		DeleteStudent:     delete_student.New(pool),
		GetStudentFinance: get_student_finance.New(pool),
		GetStudents:       get_students.New(pool),
		SearchStudent:     search_student.New(pool),
		GetStudentByID:    student_by_id.New(pool),
		MoveStudents:      move_students.New(pool),

		// Финансы
		GetAllFinance: get_all_finance.New(pool),

		// Предметы
		GetAllSubjects: get_all_subjects.New(pool),

		// AUTH
		Auth: auth.NewAggregator(pool, smtp, config),

		// Админы
		CreateAdmin:  create_admin.New(pool),
		DeleteAdmin:  delete_admin.New(pool),
		GetAdmins:    get_admins.New(pool),
		GetAdminByID: get_admin_by_id.New(pool),
	}
}
