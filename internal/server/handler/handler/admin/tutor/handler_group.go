package tutor

import (
	"github.com/it-chep/tutors.git/internal/module/admin"
	"github.com/it-chep/tutors.git/internal/server/handler/handler/admin/tutor/archive_tutor"
	"github.com/it-chep/tutors.git/internal/server/handler/handler/admin/tutor/conduct_lesson"
	"github.com/it-chep/tutors.git/internal/server/handler/handler/admin/tutor/conduct_trial"
	"github.com/it-chep/tutors.git/internal/server/handler/handler/admin/tutor/create_tutor"
	"github.com/it-chep/tutors.git/internal/server/handler/handler/admin/tutor/delete_tutor"
	"github.com/it-chep/tutors.git/internal/server/handler/handler/admin/tutor/filter_tutors"
	tutor_get_archive "github.com/it-chep/tutors.git/internal/server/handler/handler/admin/tuto
	"github.com/it-chep/tutors.git/internal/server/handler/handler/admin/tutor/get_lessons"
	"github.com/it-chep/tutors.git/internal/server/handler/handler/admin/tutor/get_tutor_finance"
	"github.com/it-chep/tutors.git/internal/server/handler/handler/admin/tutor/get_tutors"
	"github.com/it-chep/tutors.git/internal/server/handler/handler/admin/tutor/search_tutor"
	"github.com/it-chep/tutors.git/internal/server/handler/handler/admin/tutor/tutor_by_id"
	"github.com/it-chep/tutors.git/internal/server/handler/handler/admin/tutor/unarchivate_tutor"
	"github.com/it-chep/tutors.git/internal/server/handler/handler/admin/tutor/update_tutor"
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
	GetLessons    *get_lessons.Handler

	ArchivateTutor   *archive_tutor.Handler
	UnArchivateTutor *unarchivate_tutor.Handler
	UpdateTutor      *update_tutor.Handler
	FilterTutors     *filter_tutors.Handler
	GetArchive       *tutor_get_archive.Handler
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
		GetLessons:    get_lessons.NewHandler(adminModule),

		ArchivateTutor:   archive_tutor.NewHandler(adminModule),
		UnArchivateTutor: unarchivate_tutor.NewHandler(adminModule),
		UpdateTutor:      update_tutor.NewHandler(adminModule),
		FilterTutors:     filter_tutors.NewHandler(adminModule),
		GetArchive:       tutor_get_archive.NewHandler(adminModule),
	}
}
