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

func (h *Handler) setupRoutes() {}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.router.ServeHTTP(w, r)
}
