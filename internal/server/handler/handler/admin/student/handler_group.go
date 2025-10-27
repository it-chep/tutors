package student

import (
	"github.com/it-chep/tutors.git/internal/module/admin"
	"github.com/it-chep/tutors.git/internal/server/handler/handler/admin/student/create_student"
	"github.com/it-chep/tutors.git/internal/server/handler/handler/admin/student/delete_student"
	"github.com/it-chep/tutors.git/internal/server/handler/handler/admin/student/get_lessons"
	"github.com/it-chep/tutors.git/internal/server/handler/handler/admin/student/get_student_finance"
	"github.com/it-chep/tutors.git/internal/server/handler/handler/admin/student/get_students"
	"github.com/it-chep/tutors.git/internal/server/handler/handler/admin/student/move_student"
	"github.com/it-chep/tutors.git/internal/server/handler/handler/admin/student/search_student"
	"github.com/it-chep/tutors.git/internal/server/handler/handler/admin/student/student_by_id"
	"github.com/it-chep/tutors.git/internal/server/handler/handler/admin/student/update_student"
	"github.com/it-chep/tutors.git/internal/server/handler/handler/admin/student/update_wallet"
)

type HandlerGroup struct {
	GetStudentByID *student_by_id.Handler
	SearchStudent  *search_student.Handler
	GetStudents    *get_students.Handler

	CreateStudent *create_student.Handler
	DeleteStudent *delete_student.Handler

	GetStudentFinance *get_student_finance.Handler
	MoveStudent       *move_student.Handler

	UpdateWallet  *update_wallet.Handler
	GetLessons    *get_lessons.Handler
	UpdateStudent *update_student.Handler
}

func NewGroup(adminModule *admin.Module) *HandlerGroup {
	return &HandlerGroup{
		GetStudentByID: student_by_id.NewHandler(adminModule),
		SearchStudent:  search_student.NewHandler(adminModule),
		GetStudents:    get_students.NewHandler(adminModule),

		CreateStudent: create_student.NewHandler(adminModule),
		DeleteStudent: delete_student.NewHandler(adminModule),

		GetStudentFinance: get_student_finance.NewHandler(adminModule),
		MoveStudent:       move_student.NewHandler(adminModule),

		UpdateWallet:  update_wallet.NewHandler(adminModule),
		GetLessons:    get_lessons.NewHandler(adminModule),
		UpdateStudent: update_student.NewHandler(adminModule),
	}
}
