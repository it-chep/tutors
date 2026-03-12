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
	h.router.Post("/payment/{hash}", h.adminAgg.GeneratePaymentURL.Handle()) // POST /payment/{hash}

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
		r.Use(h.adminAgg.Audit.Middleware())

		r.Get("/user", h.adminAgg.Auth.GetUserInfo.Handle()) // GET /admin/user

		// Админы
		r.Route("/admins", func(r chi.Router) {
			r.Get("/", h.adminAgg.Admins.GetAdmins.Handle())                // GET /admin/admins
			r.Post("/", h.adminAgg.Admins.CreateAdmin.Handle())             // POST /admin/admins
			r.Get("/{admin_id}", h.adminAgg.Admins.GetAdminByID.Handle())   // GET /admin/admins/{id}
			r.Delete("/{admin_id}", h.adminAgg.Admins.DeleteAdmin.Handle()) // DELETE /admin/admins/{id}
		})

		// Ассистенты
		r.Route("/assistant", func(r chi.Router) {
			r.Get("/", h.adminAgg.Assistant.GetAssistants.Handle())                                        // GET /admin/assistant
			r.Post("/", h.adminAgg.Assistant.CreateAssistant.Handle())                                     // POST /admin/assistant
			r.Get("/{assistant_id}", h.adminAgg.Assistant.GetAssistantByID.Handle())                       // GET /admin/assistant/{id}
			r.Delete("/{assistant_id}", h.adminAgg.Assistant.DeleteAssistant.Handle())                     // DELETE /admin/assistant/{id}
			r.Post("/{assistant_id}/permissions", h.adminAgg.Assistant.Permissions.Handle())               // POST /admin/assistant/{id}/permissions
			r.Post("/{assistant_id}/add_available_tg", h.adminAgg.Assistant.AddAvailableTG.Handle())       // POST /admin/assistant/{id}/add_available_tg
			r.Post("/{assistant_id}/delete_available_tg", h.adminAgg.Assistant.DeleteAvailableTG.Handle()) // POST /admin/assistant/{id}/delete_available_tg
			r.Post("/{assistant_id}/penalties-bonuses", h.adminAgg.Assistant.PenaltiesBonuses.Handle())    // POST /admin/assistant/{id}/penalties-bonuses
			r.Post("/{assistant_id}/accruals", h.adminAgg.Assistant.GetAccruals.Handle())                  // POST /admin/assistant/{id}/accruals
		})

		// Репетиторы
		r.Route("/tutors", func(r chi.Router) {
			r.Get("/", h.adminAgg.Tutors.GetTutors.Handle())                             // GET /admin/tutors
			r.Get("/search", h.adminAgg.Tutors.SearchTutor.Handle())                     // GET /admin/tutors/search
			r.Get("/archive", h.adminAgg.Tutors.GetArchive.Handle())                     // GET /admin/tutors/archive
			r.Get("/contracts/download_all", h.adminAgg.Tutors.Contract.DownloadAll())   // GET /admin/tutors/contracts/download_all
			r.Post("/receipts/download_all", h.adminAgg.Tutors.Receipts.DownloadAll())   // POST /admin/tutors/receipts/download_all
			r.Post("/filter", h.adminAgg.Tutors.FilterTutors.Handle())                   // POST /admin/tutors/filter
			r.Get("/{tutor_id}", h.adminAgg.Tutors.GetTutorByID.Handle())                // GET /admin/tutors/{id}
			r.Post("/", h.adminAgg.Tutors.CreateTutor.Handle())                          // POST /admin/tutors
			r.Delete("/{tutor_id}", h.adminAgg.Tutors.DeleteTutor.Handle())              // DELETE /admin/tutors/{id}
			r.Post("/{tutor_id}/contract", h.adminAgg.Tutors.Contract.Upload())          // POST /admin/tutors/{id}/contract
			r.Get("/{tutor_id}/contract", h.adminAgg.Tutors.Contract.Download())         // GET /admin/tutors/{id}/contract
			r.Delete("/{tutor_id}/contract", h.adminAgg.Tutors.Contract.Delete())        // DELETE /admin/tutors/{id}/contract
			r.Post("/{tutor_id}/finance", h.adminAgg.Tutors.GetTutorFinance.Handle())    // POST /admin/tutors/{id}/finance
			r.Post("/trial_lesson", h.adminAgg.Tutors.ConductTrial.Handle())             // POST /admin/tutors/trial_lesson
			r.Post("/conduct_lesson", h.adminAgg.Tutors.ConductLesson.Handle())          // POST /admin/tutors/conduct_lesson
			r.Post("/{tutor_id}/lessons", h.adminAgg.Tutors.GetLessons.Handle())         // POST /admin/tutors/{id}/lessons
			r.Post("/{tutor_id}/archive", h.adminAgg.Tutors.ArchivateTutor.Handle())     // POST /admin/tutors/{id}/archive
			r.Post("/{tutor_id}/unarchive", h.adminAgg.Tutors.UnArchivateTutor.Handle()) // POST /admin/tutors/{id}/unarchive
			r.Post("/{tutor_id}/update", h.adminAgg.Tutors.UpdateTutor.Handle())         // POST /admin/tutors/{id}/update
			r.Post("/{tutor_id}/penalties-bonuses", h.adminAgg.Tutors.PenaltiesBonuses.Handle())
			r.Post("/{tutor_id}/accruals", h.adminAgg.Tutors.GetAccruals.Handle())
			r.Post("/{tutor_id}/payouts", h.adminAgg.Tutors.Payouts.Handle())
			r.Post("/{tutor_id}/receipts", h.adminAgg.Tutors.Receipts.Handle())
		})

		// Студенты
		r.Route("/students", func(r chi.Router) {
			r.Get("/", h.adminAgg.Students.GetStudents.Handle())                                           // GET /admin/students
			r.Get("/search", h.adminAgg.Students.SearchStudent.Handle())                                   // GET /admin/students/search
			r.Get("/archive", h.adminAgg.Students.GetArchive.Handle())                                     // GET /admin/students/archive
			r.Post("/push_all_students", h.adminAgg.Students.PushAllDebitors.Handle())                     // POST /admin/students/push_all_students
			r.Get("/tg_admins_usernames", h.adminAgg.Students.GetTgAdminsUsernames.Handle())               // GET /admin/students/tg_admins_usernames
			r.Post("/", h.adminAgg.Students.CreateStudent.Handle())                                        // POST /admin/students
			r.Post("/filter", h.adminAgg.Students.FilterStudents.Handle())                                 // POST /admin/students/filter
			r.Post("/move", h.adminAgg.Students.MoveStudent.Handle())                                      // POST /admin/students/move
			r.Post("/change_all_payment", h.adminAgg.Students.ChangeAllPayment.Handle())                   // POST /admin/students/change_all_payment
			r.Get("/{student_id}", h.adminAgg.Students.GetStudentByID.Handle())                            // GET /admin/students/{id}
			r.Delete("/{student_id}", h.adminAgg.Students.DeleteStudent.Handle())                          // DELETE /admin/students/{id}
			r.Post("/{student_id}", h.adminAgg.Students.UpdateStudent.Handle())                            // POST /admin/students/{id}
			r.Post("/{student_id}/finance", h.adminAgg.Students.GetStudentFinance.Handle())                // POST /admin/students/{id}/finance
			r.Post("/{student_id}/wallet", h.adminAgg.Students.UpdateWallet.Handle())                      // POST /admin/students/{id}/wallet
			r.Post("/{student_id}/change_payment", h.adminAgg.Students.ChangeStudentPayment.Handle())      // POST /admin/students/{id}/change_payment
			r.Post("/{student_id}/lessons", h.adminAgg.Students.GetLessons.Handle())                       // POST /admin/students/{id}/lessons
			r.Post("/{student_id}/transactions", h.adminAgg.Students.GetTransactionHistory.Handle())       // POST /admin/students/{id}/transactions
			r.Post("/{student_id}/transactions/manual", h.adminAgg.Students.AddManualTransaction.Handle()) // POST /admin/students/{id}/transactions/manual
			r.Post("/{student_id}/notifications", h.adminAgg.Students.GetNotificationHistory.Handle())     // POST /admin/students/{id}/notifications
			r.Post("/{student_id}/notifications/push", h.adminAgg.Students.PushNotification.Handle())      // POST /admin/students/{id}/notifications/push
			r.Post("/{student_id}/comments", h.adminAgg.Students.CreateComment.Handle())                   // POST /admin/students/{id}/comments
			r.Get("/{student_id}/comments", h.adminAgg.Students.GetComments.Handle())                      // GET /admin/students/{id}/comments
			r.Delete("/{student_id}/comments/{comment_id}", h.adminAgg.Students.DeleteComment.Handle())    // DELETE /admin/students/{id}/comments/{comment_id}
			r.Post("/{student_id}/archive", h.adminAgg.Students.ArchiveStudent.Handle())                   // POST /admin/students/{id}/archive
			r.Post("/{student_id}/unarchive", h.adminAgg.Students.UnArchivateStudent.Handle())             // POST /admin/students/{id}/unarchive
		})

		// Уроки
		r.Route("/lessons", func(r chi.Router) {
			r.Delete("/{lesson_id}", h.adminAgg.Lessons.DeleteLesson.Handle()) // DELETE /admin/lessons/{id}
			r.Post("/{lesson_id}", h.adminAgg.Lessons.UpdateLesson.Handle())   // POST /admin/lessons/{id}
			r.Post("/", h.adminAgg.GetAllLessons.Handle())                     // POST /admin/lessons
		})

		r.Get("/subjects", h.adminAgg.GetAllSubjects.Handle())            // GET /admin/subjects
		r.Post("/finance", h.adminAgg.GetAllFinance.Handle())             // POST /admin/finance
		r.Post("/finance_by_tgs", h.adminAgg.GetAllFinanceByTGs.Handle()) // POST /admin/finance_by_tgs
		r.Post("/transactions", h.adminAgg.GetAllTransactions.Handle())   // POST /admin/transactions
		r.Get("/payments", h.adminAgg.GetAdminPayments.Handle())          // POST /admin/payments
	})

	h.router.Route("/tutors", func(r chi.Router) {
		r.Use(h.adminAgg.Auth.CheckPathPermission.AuthMiddleware())
		r.Post("/save_receipt", h.adminAgg.Tutors.SaveReceipt.Handle()) // POST /tutors/save_receipt
	})

	h.router.Post("/webhook/alpha", h.adminAgg.AlphaHook.Handle())      // POST /alpha/hook
	h.router.Post("/callback/tbank", h.adminAgg.TbankCallBack.Handle()) // POST /callback/tbank
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.router.ServeHTTP(w, r)
}
