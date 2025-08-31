package server

import (
	"github.com/it-chep/tutors.git/internal/server/handler"
	"net/http"
	"time"
)

type Server struct {
	srv *http.Server
}

func New(handler *handler.Handler) *Server {
	srv := &http.Server{
		// todo: подумать
		Addr:         ":8080",
		Handler:      handler,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	return &Server{srv: srv}
}

func (s *Server) ListenAndServe() error {
	return s.srv.ListenAndServe()
}
