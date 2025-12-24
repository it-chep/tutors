package action

import (
	"github.com/it-chep/tutors.git/internal/config"
	"github.com/it-chep/tutors.git/internal/module/admin/action/admin/create_admin"
	"github.com/it-chep/tutors.git/internal/module/admin/action/admin/delete_admin"
	"github.com/it-chep/tutors.git/internal/module/admin/action/admin/get_admin_by_id"
	"github.com/it-chep/tutors.git/internal/module/admin/action/admin/get_admins"
	"github.com/it-chep/tutors.git/internal/module/admin/action/assistant/add_available_tg"
	"github.com/it-chep/tutors.git/internal/module/admin/action/assistant/delete_available_tg"
	"github.com/it-chep/tutors.git/internal/module/admin/action/assistant/get_assistance"
	"github.com/it-chep/tutors.git/internal/module/admin/action/assistant/get_available_tg"
	"github.com/it-chep/tutors.git/internal/module/admin/action/auth"
	"github.com/it-chep/tutors.git/internal/module/admin/action/get_all_finance"
	"github.com/it-chep/tutors.git/internal/module/admin/action/get_all_finance_by_tgs"
	"github.com/it-chep/tutors.git/internal/module/admin/action/get_all_lessons"
	"github.com/it-chep/tutors.git/internal/module/admin/action/get_all_subjects"
	"github.com/it-chep/tutors.git/internal/module/admin/action/get_all_transactions"
	"github.com/it-chep/tutors.git/internal/module/admin/action/lessons/delete_lesson"
	"github.com/it-chep/tutors.git/internal/module/admin/action/lessons/update_lesson"
	"github.com/it-chep/tutors.git/internal/module/admin/action/student/archivate_student"
	"github.com/it-chep/tutors.git/internal/module/admin/action/student/archive_filter"
	"github.com/it-chep/tutors.git/internal/module/admin/action/student/filter_students"
	"github.com/it-chep/tutors.git/internal/module/admin/action/student/get_archive"
	"github.com/it-chep/tutors.git/internal/module/admin/action/student/get_notification_history"
	"github.com/it-chep/tutors.git/internal/module/admin/action/student/get_tg_admins_usernames"
	"github.com/it-chep/tutors.git/internal/module/admin/action/student/get_transaction_history"
	"github.com/it-chep/tutors.git/internal/module/admin/action/student/push_all_debitors"
	"github.com/it-chep/tutors.git/internal/module/admin/action/student/push_notification"
	"github.com/it-chep/tutors.git/internal/module/admin/action/student/unarchivate_student"

	"github.com/it-chep/tutors.git/internal/module/admin/action/student/create_student"
	"github.com/it-chep/tutors.git/internal/module/admin/action/student/delete_student"
	"github.com/it-chep/tutors.git/internal/module/admin/action/student/get_student_finance"
	"github.com/it-chep/tutors.git/internal/module/admin/action/student/get_student_lessons"
	"github.com/it-chep/tutors.git/internal/module/admin/action/student/get_students"
	"github.com/it-chep/tutors.git/internal/module/admin/action/student/move_students"
	"github.com/it-chep/tutors.git/internal/module/admin/action/student/search_student"
	"github.com/it-chep/tutors.git/internal/module/admin/action/student/student_by_id"
	"github.com/it-chep/tutors.git/internal/module/admin/action/student/update_student"
	"github.com/it-chep/tutors.git/internal/module/admin/action/student/update_wallet"
	"github.com/it-chep/tutors.git/internal/module/admin/action/tutor/conduct_lesson"
	"github.com/it-chep/tutors.git/internal/module/admin/action/tutor/conduct_trial"
	"github.com/it-chep/tutors.git/internal/module/admin/action/tutor/create_tutor"
	"github.com/it-chep/tutors.git/internal/module/admin/action/tutor/delete_tutor"
	"github.com/it-chep/tutors.git/internal/module/admin/action/tutor/get_tutor_finance"
	"github.com/it-chep/tutors.git/internal/module/admin/action/tutor/get_tutor_lessons"
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
	GetTutorLessons *get_tutor_lessons.Action

	// Студент
	CreateStudent          *create_student.Action
	DeleteStudent          *delete_student.Action
	GetStudentFinance      *get_student_finance.Action
	GetStudents            *get_students.Action
	FilterStudents         *filter_students.Action
	SearchStudent          *search_student.Action
	GetStudentByID         *student_by_id.Action
	MoveStudents           *move_students.Action
	UpdateWallet           *update_wallet.Action
	UpdateStudent          *update_student.Action
	GetStudentLessons      *get_student_lessons.Action
	GetTgAdminsUsernames   *get_tg_admins_usernames.Action
	GetTransactionHistory  *get_transaction_history.Action
	GetNotificationHistory *get_notification_history.Action
	// пуши
	PushNotification *push_notification.Action
	PushAllDebitors  *push_all_debitors.Action
	// архив
	GetArchive        *get_archive.Action
	ArchiveStudent    *archivate_student.Action
	UnArhivateStudent *unarchivate_student.Action
	ArchiveFilter     *archive_filter.Action

	// Финансы
	GetAllFinance      *get_all_finance.Action
	GetAllFinanceByTGs *get_all_finance_by_tgs.Action

	// Предметы
	GetAllSubjects *get_all_subjects.Action

	// Транзакции
	GetAllTransactions *get_all_transactions.Action

	// Уроки
	GetAllLessons *get_all_lessons.Action

	// AUTH
	Auth *auth.Aggregator

	// Админы
	CreateAdmin  *create_admin.Action
	DeleteAdmin  *delete_admin.Action
	GetAdmins    *get_admins.Action
	GetAdminByID *get_admin_by_id.Action

	// Уроки
	DeleteLesson *delete_lesson.Action
	UpdateLesson *update_lesson.Action

	// Ассистенты
	AddAvailableTg           *add_available_tg.Action
	DeleteAvailableTg        *delete_available_tg.Action
	GetAssistantAvailableTGs *get_available_tg.Action
	GetAssistance            *get_assistance.Action
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
		GetTutorLessons: get_tutor_lessons.New(pool),

		// Студент
		CreateStudent:          create_student.New(pool),
		DeleteStudent:          delete_student.New(pool),
		GetStudentFinance:      get_student_finance.New(pool),
		GetStudents:            get_students.New(pool),
		FilterStudents:         filter_students.New(pool),
		SearchStudent:          search_student.New(pool),
		GetStudentByID:         student_by_id.New(pool),
		MoveStudents:           move_students.New(pool),
		UpdateWallet:           update_wallet.New(pool),
		UpdateStudent:          update_student.New(pool),
		GetStudentLessons:      get_student_lessons.New(pool),
		GetTgAdminsUsernames:   get_tg_admins_usernames.New(pool),
		GetTransactionHistory:  get_transaction_history.New(pool),
		GetNotificationHistory: get_notification_history.New(pool),
		// пуши
		PushNotification: push_notification.New(pool, bot),
		PushAllDebitors:  push_all_debitors.New(pool, bot),
		// архив
		GetArchive:        get_archive.New(pool),
		ArchiveStudent:    archivate_student.New(pool),
		UnArhivateStudent: unarchivate_student.New(pool),
		ArchiveFilter:     archive_filter.New(pool),

		// Финансы
		GetAllFinance:      get_all_finance.New(pool),
		GetAllFinanceByTGs: get_all_finance_by_tgs.New(pool),

		// Предметы
		GetAllSubjects: get_all_subjects.New(pool),

		// Транзакции
		GetAllTransactions: get_all_transactions.New(pool),

		// Уроки
		GetAllLessons: get_all_lessons.New(pool),

		// AUTH
		Auth: auth.NewAggregator(pool, smtp, config),

		// Админы
		CreateAdmin:  create_admin.New(pool),
		DeleteAdmin:  delete_admin.New(pool),
		GetAdmins:    get_admins.New(pool),
		GetAdminByID: get_admin_by_id.New(pool),

		// Уроки
		DeleteLesson: delete_lesson.New(pool),
		UpdateLesson: update_lesson.New(pool, bot),

		// Ассистенты
		AddAvailableTg:           add_available_tg.New(pool),
		DeleteAvailableTg:        delete_available_tg.New(pool),
		GetAssistantAvailableTGs: get_available_tg.New(pool),
		GetAssistance:            get_assistance.New(pool),
	}
}
