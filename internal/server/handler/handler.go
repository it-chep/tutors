package handler

import (
	"fmt"
	"net/http"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/it-chep/tutors.git/internal/module/bot"

	"github.com/go-chi/chi/v5"
	"github.com/it-chep/tutors.git/internal/module/admin"
	adminHandler "github.com/it-chep/tutors.git/internal/server/handler/handler/admin"
	"github.com/it-chep/tutors.git/internal/server/middleware"
)

type Config interface {
	Token() string
}

// TgHookParser .
type TgHookParser interface {
	HandleUpdate(r *http.Request) (*tgbotapi.Update, error)
}

type Handler struct {
	router    *chi.Mux
	botParser TgHookParser

	botModule *bot.Bot

	adminAgg *adminHandler.HandlerAggregator
}

func NewHandler(botParser TgHookParser, botModule *bot.Bot, adminModule *admin.Module, cfg Config) *Handler {
	h := &Handler{
		router:    chi.NewRouter(),
		botParser: botParser,
		botModule: botModule,
	}

	h.setupMiddleware()
	h.setupHandlerAggregator(adminModule)
	h.setupRoutes(cfg)

	return h
}

func (h *Handler) setupMiddleware() {
	h.router.Use(middleware.LoggerMiddleware)
	h.router.Use(middleware.CORSMiddleware)
}

func (h *Handler) setupHandlerAggregator(adminModule *admin.Module) {
	h.adminAgg = adminHandler.NewAggregator(adminModule)
}

func (h *Handler) setupRoutes(cfg Config) {
	h.router.Post(fmt.Sprintf("/%s/", cfg.Token()), h.bot())

	h.router.Route("/auth", func(r chi.Router) {
		// Регистрация
		r.Route("/register", func(r chi.Router) {
			r.Post("/", h.adminAgg.Auth.Register.RegisterHandler())     // POST /auth/register
			r.Post("/verify", h.adminAgg.Auth.Register.VerifyHandler()) // POST /auth/register/verify
		})

		// Авторизация
		r.Route("/login", func(r chi.Router) {
			r.Post("/", h.adminAgg.Auth.Login.LoginHandler())        // POST /auth/login
			r.Post("/verify", h.adminAgg.Auth.Login.VerifyHandler()) // POST /auth/login/verify
		})

		// Аутентификация
		r.Route("/auth", func(r chi.Router) {
			r.Get("/refresh", h.adminAgg.Auth.Refresh.RefreshHandler()) // GET /auth/refresh
		})

		// Выход
		r.Route("/logout", func(r chi.Router) {
			r.Post("/", h.adminAgg.Auth.Logout.DeleteCookieHandler())
		})
	})

	h.router.Route("/admin", func(r chi.Router) {
		r.Use(h.adminAgg.Auth.CheckPathPermission.AuthMiddleware())

		r.Get("/user", h.adminAgg.Auth.GetUserInfo.Handle()) // GET /admin/user

		// Админы
		r.Route("/admins", func(r chi.Router) {
			r.Get("/", h.adminAgg.Admins.GetAdmins.Handle())                // GET /admin/admins
			r.Post("/", h.adminAgg.Admins.CreateAdmin.Handle())             // POST /admin/admins
			r.Get("/{admin_id}", h.adminAgg.Admins.GetAdminByID.Handle())   // GET /admin/admins/{id}
			r.Delete("/{admin_id}", h.adminAgg.Admins.DeleteAdmin.Handle()) // DELETE /admin/admins/{id}
		})

		// Репетиторы
		r.Route("/tutors", func(r chi.Router) {
			r.Get("/", h.adminAgg.Tutors.GetTutors.Handle())                          // GET /admin/tutors
			r.Get("/search", h.adminAgg.Tutors.SearchTutor.Handle())                  // GET /admin/tutors/search
			r.Get("/{tutor_id}", h.adminAgg.Tutors.GetTutorByID.Handle())             // GET /admin/tutors/{id}
			r.Post("/", h.adminAgg.Tutors.CreateTutor.Handle())                       // POST /admin/tutors
			r.Delete("/{tutor_id}", h.adminAgg.Tutors.DeleteTutor.Handle())           // DELETE /admin/tutors/{id}
			r.Post("/{tutor_id}/finance", h.adminAgg.Tutors.GetTutorFinance.Handle()) // POST /admin/tutors/{id}/finance
			r.Post("/trial_lesson", h.adminAgg.Tutors.ConductTrial.Handle())          // POST /admin/tutors/trial_lesson
			r.Post("/conduct_lesson", h.adminAgg.Tutors.ConductLesson.Handle())       // POST /admin/tutors/conduct_lesson
		})

		// Студенты
		r.Route("/students", func(r chi.Router) {
			r.Get("/", h.adminAgg.Students.GetStudents.Handle())                            // GET /admin/students
			r.Get("/search", h.adminAgg.Students.SearchStudent.Handle())                    // GET /admin/students/search
			r.Get("/{student_id}", h.adminAgg.Students.GetStudentByID.Handle())             // GET /admin/students/{id}
			r.Post("/", h.adminAgg.Students.CreateStudent.Handle())                         // POST /admin/students
			r.Delete("/{student_id}", h.adminAgg.Students.DeleteStudent.Handle())           // DELETE /admin/students/{id}
			r.Post("/{student_id}/finance", h.adminAgg.Students.GetStudentFinance.Handle()) // POST /admin/students/{id}/finance
			r.Post("/move", h.adminAgg.Students.MoveStudent.Handle())                       // POST /admin/students/move
		})

		r.Get("/subjects", h.adminAgg.GetAllSubjects.Handle()) // GET /admin/subjects
		r.Post("/finance", h.adminAgg.GetAllFinance.Handle())  // POST /admin/finance
	})

	h.router.Post("/webhook/alpha", h.adminAgg.AlphaHook.Handle())      // POST /alpha/hook
	h.router.Post("/callback/tbank", h.adminAgg.TbankCallBack.Handle()) // POST /callback/tbank
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.router.ServeHTTP(w, r)
}
