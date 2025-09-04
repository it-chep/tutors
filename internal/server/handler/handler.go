package handler

import (
	"github.com/go-chi/chi/v5"
	"github.com/it-chep/tutors.git/internal/module/admin"
	adminHandler "github.com/it-chep/tutors.git/internal/server/handler/handler/admin"
	"github.com/it-chep/tutors.git/internal/server/middleware"
	"net/http"
)

type Handler struct {
	router *chi.Mux

	adminAgg *adminHandler.HandlerAggregator
}

func NewHandler(adminModule *admin.Module) *Handler {
	h := &Handler{
		router: chi.NewRouter(),
	}

	//h.setupMiddleware()
	h.setupHandlerAggregator(adminModule)
	//h.setupRoutes(cfg)

	return h
}

func (h *Handler) setupMiddleware() {
	h.router.Use(middleware.LoggerMiddleware)
}

func (h *Handler) setupHandlerAggregator(adminModule *admin.Module) {
	h.adminAgg = adminHandler.NewAggregator(adminModule)
}

func (h *Handler) setupRoutes() {
	h.router.Route("/", func(r chi.Router) {
		//r.Post(fmt.Sprintf("/%s/", cfg.Token()), h.bot())
	})

	h.router.Route("/admin", func(r chi.Router) {
		//r.Get("/", h.admin())
		r.Get("/roles", h.adminAgg.GetAvailableRoles.Handle()) // GET /admin/roles

		// Админы
		r.Route("/admins", func(r chi.Router) {

		})

		// Репетиторы
		r.Route("/tutors", func(r chi.Router) {
			r.Get("/", h.adminAgg.Tutors.GetTutors.Handle())                            // GET /admin/tutors
			r.Get("/search", h.adminAgg.Tutors.SearchTutor.Handle())                    // GET /admin/tutors/search
			r.Get("/{student_id}", h.adminAgg.Tutors.GetTutorByID.Handle())             // GET /admin/tutors/{id}
			r.Post("/", h.adminAgg.Tutors.CreateTutor.Handle())                         // POST /admin/tutors
			r.Delete("/{student_id}", h.adminAgg.Tutors.DeleteTutor.Handle())           // DELETE /admin/tutors/{id}
			r.Post("/{student_id}/finance", h.adminAgg.Tutors.GetTutorFinance.Handle()) // POST /admin/tutors/{id}/finance
			r.Post("/trial_lesson", h.adminAgg.Tutors.ConductTrial.Handle())            // POST /admin/tutors/trial_lesson
			r.Post("/conduct_lesson", h.adminAgg.Tutors.ConductLesson.Handle())         // POST /admin/tutors/conduct_lesson
		})

		// Студенты
		r.Route("/students", func(r chi.Router) {
			r.Get("/", h.adminAgg.Students.GetStudents.Handle())                            // GET /admin/students
			r.Get("/search", h.adminAgg.Students.SearchStudent.Handle())                    // GET /admin/students/search
			r.Get("/{student_id}", h.adminAgg.Students.GetStudentByID.Handle())             // GET /admin/students/{id}
			r.Post("/", h.adminAgg.Students.CreateStudent.Handle())                         // POST /admin/students
			r.Delete("/{student_id}", h.adminAgg.Students.DeleteStudent.Handle())           // DELETE /admin/students/{id}
			r.Post("/{student_id}/finance", h.adminAgg.Students.GetStudentFinance.Handle()) // POST /admin/students/{id}/finance
		})

		r.Get("/subjects", h.adminAgg.GetAllSubjects.Handle()) // GET /admin/subjects
		r.Post("/finance", h.adminAgg.GetAllFinance.Handle())  // POST /admin/finance
	})

}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.router.ServeHTTP(w, r)
}
