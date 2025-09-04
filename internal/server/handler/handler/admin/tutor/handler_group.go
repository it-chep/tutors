package tutor

import (
	"github.com/it-chep/tutors.git/internal/module/admin"
	"github.com/it-chep/tutors.git/internal/server/handler/handler/admin/tutor/conduct_lesson"
	"github.com/it-chep/tutors.git/internal/server/handler/handler/admin/tutor/conduct_trial"
	"github.com/it-chep/tutors.git/internal/server/handler/handler/admin/tutor/create_tutor"
	"github.com/it-chep/tutors.git/internal/server/handler/handler/admin/tutor/delete_tutor"
	"github.com/it-chep/tutors.git/internal/server/handler/handler/admin/tutor/get_tutor_finance"
	"github.com/it-chep/tutors.git/internal/server/handler/handler/admin/tutor/get_tutors"
	"github.com/it-chep/tutors.git/internal/server/handler/handler/admin/tutor/search_tutor"
	"github.com/it-chep/tutors.git/internal/server/handler/handler/admin/tutor/tutor_by_id"
)

type HandlerGroup struct {
	CreateTutor *create_tutor.Handler
	DeleteTutor *delete_tutor.Handler

	GetTutorFinance *get_tutor_finance.Handler
	GetTutorByID    *tutor_by_id.Handler
	GetTutors       *get_tutors.Handler

	SearchTutor *search_tutor.Handler

	ConductTrial  *conduct_trial.Handler
	ConductLesson *conduct_lesson.Handler
}

func NewGroup(adminModule *admin.Module) *HandlerGroup {
	return &HandlerGroup{
		CreateTutor: create_tutor.NewHandler(adminModule),
		DeleteTutor: delete_tutor.NewHandler(adminModule),

		GetTutorFinance: get_tutor_finance.NewHandler(adminModule),
		GetTutorByID:    tutor_by_id.NewHandler(adminModule),
		GetTutors:       get_tutors.NewHandler(adminModule),

		SearchTutor: search_tutor.NewHandler(adminModule),

		ConductTrial:  conduct_trial.NewHandler(adminModule),
		ConductLesson: conduct_lesson.NewHandler(adminModule),
	}
}
