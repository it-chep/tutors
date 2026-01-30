package student

import (
	"github.com/it-chep/tutors.git/internal/module/admin"
	"github.com/it-chep/tutors.git/internal/server/handler/handler/admin/student/archive_student"
	"github.com/it-chep/tutors.git/internal/server/handler/handler/admin/student/change_student_payment"
	"github.com/it-chep/tutors.git/internal/server/handler/handler/admin/student/create_student"
	"github.com/it-chep/tutors.git/internal/server/handler/handler/admin/student/delete_student"
	"github.com/it-chep/tutors.git/internal/server/handler/handler/admin/student/filter_students"
	"github.com/it-chep/tutors.git/internal/server/handler/handler/admin/student/get_archive"
	"github.com/it-chep/tutors.git/internal/server/handler/handler/admin/student/get_lessons"
	"github.com/it-chep/tutors.git/internal/server/handler/handler/admin/student/get_notification_history"
	"github.com/it-chep/tutors.git/internal/server/handler/handler/admin/student/get_student_finance"
	"github.com/it-chep/tutors.git/internal/server/handler/handler/admin/student/get_students"
	"github.com/it-chep/tutors.git/internal/server/handler/handler/admin/student/get_tg_admins_usernames"
	"github.com/it-chep/tutors.git/internal/server/handler/handler/admin/student/get_transaction_history"
	"github.com/it-chep/tutors.git/internal/server/handler/handler/admin/student/move_student"
	"github.com/it-chep/tutors.git/internal/server/handler/handler/admin/student/push_all_debitors"
	"github.com/it-chep/tutors.git/internal/server/handler/handler/admin/student/push_notification"
	"github.com/it-chep/tutors.git/internal/server/handler/handler/admin/student/search_student"
	"github.com/it-chep/tutors.git/internal/server/handler/handler/admin/student/student_by_id"
	"github.com/it-chep/tutors.git/internal/server/handler/handler/admin/student/unarchivate_student"
	"github.com/it-chep/tutors.git/internal/server/handler/handler/admin/student/update_student"
	"github.com/it-chep/tutors.git/internal/server/handler/handler/admin/student/update_wallet"
)

type HandlerGroup struct {
	GetStudentByID *student_by_id.Handler
	SearchStudent  *search_student.Handler
	GetStudents    *get_students.Handler
	FilterStudents *filter_students.Handler

	CreateStudent *create_student.Handler
	DeleteStudent *delete_student.Handler

	GetStudentFinance *get_student_finance.Handler
	MoveStudent       *move_student.Handler

	UpdateWallet           *update_wallet.Handler
	GetLessons             *get_lessons.Handler
	UpdateStudent          *update_student.Handler
	GetTgAdminsUsernames   *get_tg_admins_usernames.Handler
	GetTransactionHistory  *get_transaction_history.Handler
	GetNotificationHistory *get_notification_history.Handler
	PushNotification       *push_notification.Handler
	PushAllDebitors        *push_all_debitors.Handler
	ChangeStudentPayment   *change_student_payment.Handler

	GetArchive         *get_archive.Handler
	ArchiveStudent     *archive_student.Handler
	UnArchivateStudent *unarchivate_student.Handler
}

func NewGroup(adminModule *admin.Module) *HandlerGroup {
	return &HandlerGroup{
		GetStudentByID: student_by_id.NewHandler(adminModule),
		SearchStudent:  search_student.NewHandler(adminModule),
		GetStudents:    get_students.NewHandler(adminModule),
		FilterStudents: filter_students.NewHandler(adminModule),

		CreateStudent: create_student.NewHandler(adminModule),
		DeleteStudent: delete_student.NewHandler(adminModule),

		GetStudentFinance: get_student_finance.NewHandler(adminModule),
		MoveStudent:       move_student.NewHandler(adminModule),

		UpdateWallet:           update_wallet.NewHandler(adminModule),
		GetLessons:             get_lessons.NewHandler(adminModule),
		UpdateStudent:          update_student.NewHandler(adminModule),
		GetTgAdminsUsernames:   get_tg_admins_usernames.NewHandler(adminModule),
		GetTransactionHistory:  get_transaction_history.NewHandler(adminModule),
		GetNotificationHistory: get_notification_history.NewHandler(adminModule),
		PushNotification:       push_notification.NewHandler(adminModule),
		PushAllDebitors:        push_all_debitors.NewHandler(adminModule),
		ChangeStudentPayment:   change_student_payment.NewHandler(adminModule),

		GetArchive:         get_archive.NewHandler(adminModule),
		ArchiveStudent:     archive_student.NewHandler(adminModule),
		UnArchivateStudent: unarchivate_student.NewHandler(adminModule),
	}
}
