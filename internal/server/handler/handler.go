package handler

import (
	"github.com/go-chi/chi/v5"
	"github.com/it-chep/tutors.git/internal/server/middleware"
	"net/http"
)

type Handler struct {
	router *chi.Mux
}

func NewHandler() *Handler {
	h := &Handler{
		router: chi.NewRouter(),
	}

	//h.setupMiddleware()
	//h.setupHandlerAggregator(adminModule)
	//h.setupRoutes(cfg)

	return h
}

func (h *Handler) setupMiddleware() {
	h.router.Use(middleware.LoggerMiddleware)
}

func (h *Handler) setupHandlerAggregator() {
	//h.adminAgg = adminHandler.NewAggregator(adminModule)
}

func (h *Handler) setupRoutes() {
	h.router.Route("/", func(r chi.Router) {
		//r.Post(fmt.Sprintf("/%s/", cfg.Token()), h.bot())
	})

	h.router.Route("/admin", func(r chi.Router) {
		//r.Get("/", h.admin())

		// Админы
		r.Route("/admins", func(r chi.Router) {

		})
		// Репетиторы
		r.Route("/tutors", func(r chi.Router) {

		})

		// Студенты
		r.Route("/students", func(r chi.Router) {
			//r.Get("/", h.)
		})
	})

}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.router.ServeHTTP(w, r)
}
