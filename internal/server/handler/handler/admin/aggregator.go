package admin

import (
	"github.com/it-chep/tutors.git/internal/module/admin"
	"github.com/it-chep/tutors.git/internal/module/admin/action/auth"
	"github.com/it-chep/tutors.git/internal/module/admin/alpha"
	"github.com/it-chep/tutors.git/internal/module/admin/tbank"
	"github.com/it-chep/tutors.git/internal/server/handler/handler/admin/admins"
	"github.com/it-chep/tutors.git/internal/server/handler/handler/admin/get_all_finance"
	"github.com/it-chep/tutors.git/internal/server/handler/handler/admin/get_all_subjects"
	"github.com/it-chep/tutors.git/internal/server/handler/handler/admin/lessons"
	"github.com/it-chep/tutors.git/internal/server/handler/handler/admin/student"
	"github.com/it-chep/tutors.git/internal/server/handler/handler/admin/tutor"
)

type HandlerAggregator struct {
	Students *student.HandlerGroup
	Tutors   *tutor.HandlerGroup
	Admins   *admins.HandlerGroup
	Lessons  *lessons.HandlerGroup

	GetAllFinance  *get_all_finance.Handler
	GetAllSubjects *get_all_subjects.Handler

	AlphaHook     *alpha.WebHookAlpha
	TbankCallBack *tbank.CallbackTbank

	Auth *auth.Aggregator
}

func NewAggregator(adminModule *admin.Module) *HandlerAggregator {
	return &HandlerAggregator{
		Students: student.NewGroup(adminModule),
		Tutors:   tutor.NewGroup(adminModule),
		Admins:   admins.NewGroup(adminModule),
		Lessons:  lessons.NewGroup(adminModule),

		GetAllFinance:  get_all_finance.NewHandler(adminModule),
		GetAllSubjects: get_all_subjects.NewHandler(adminModule),

		AlphaHook:     adminModule.AlphaHook,
		TbankCallBack: adminModule.TbankCallback,
		Auth:          adminModule.Actions.Auth,
	}
}
