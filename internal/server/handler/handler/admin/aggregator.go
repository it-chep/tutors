package admin

import (
	"github.com/it-chep/tutors.git/internal/module/admin"
	"github.com/it-chep/tutors.git/internal/module/admin/action/auth"
	"github.com/it-chep/tutors.git/internal/module/admin/alpha"
	"github.com/it-chep/tutors.git/internal/module/admin/tbank"
	"github.com/it-chep/tutors.git/internal/server/handler/handler/admin/admins"
	"github.com/it-chep/tutors.git/internal/server/handler/handler/admin/get_all_finance"
	"github.com/it-chep/tutors.git/internal/server/handler/handler/admin/get_all_finance_by_tgs"
	"github.com/it-chep/tutors.git/internal/server/handler/handler/admin/get_all_lessons"
	"github.com/it-chep/tutors.git/internal/server/handler/handler/admin/get_all_subjects"
	"github.com/it-chep/tutors.git/internal/server/handler/handler/admin/get_all_transactions"
	"github.com/it-chep/tutors.git/internal/server/handler/handler/admin/lessons"
	"github.com/it-chep/tutors.git/internal/server/handler/handler/admin/student"
	"github.com/it-chep/tutors.git/internal/server/handler/handler/admin/tutor"
)

type HandlerAggregator struct {
	Students *student.HandlerGroup
	Tutors   *tutor.HandlerGroup
	Admins   *admins.HandlerGroup
	Lessons  *lessons.HandlerGroup

	GetAllFinance      *get_all_finance.Handler
	GetAllFinanceByTGs *get_all_finance_by_tgs.Handler
	GetAllSubjects     *get_all_subjects.Handler
	GetAllTransactions *get_all_transactions.Handler
	GetAllLessons      *get_all_lessons.Handler

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

		GetAllFinance:      get_all_finance.NewHandler(adminModule),
		GetAllFinanceByTGs: get_all_finance_by_tgs.NewHandler(adminModule),
		GetAllSubjects:     get_all_subjects.NewHandler(adminModule),
		GetAllTransactions: get_all_transactions.NewHandler(adminModule),
		GetAllLessons:      get_all_lessons.NewHandler(adminModule),

		AlphaHook:     adminModule.AlphaHook,
		TbankCallBack: adminModule.TbankCallback,
		Auth:          adminModule.Actions.Auth,
	}
}
