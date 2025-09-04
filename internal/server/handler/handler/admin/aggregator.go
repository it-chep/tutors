package admin

import (
	"github.com/it-chep/tutors.git/internal/module/admin"
	"github.com/it-chep/tutors.git/internal/server/handler/handler/admin/get_all_finance"
	"github.com/it-chep/tutors.git/internal/server/handler/handler/admin/get_all_subjects"
	"github.com/it-chep/tutors.git/internal/server/handler/handler/admin/get_available_roles"
	"github.com/it-chep/tutors.git/internal/server/handler/handler/admin/student"
	"github.com/it-chep/tutors.git/internal/server/handler/handler/admin/tutor"
)

type HandlerAggregator struct {
	Students *student.HandlerGroup
	Tutors   *tutor.HandlerGroup

	GetAllFinance     *get_all_finance.Handler
	GetAllSubjects    *get_all_subjects.Handler
	GetAvailableRoles *get_available_roles.Handler
}

func NewAggregator(adminModule *admin.Module) *HandlerAggregator {
	return &HandlerAggregator{
		Students: student.NewGroup(adminModule),
		Tutors:   tutor.NewGroup(adminModule),

		GetAllFinance:     get_all_finance.NewHandler(adminModule),
		GetAllSubjects:    get_all_subjects.NewHandler(adminModule),
		GetAvailableRoles: get_available_roles.NewHandler(adminModule),
	}
}
