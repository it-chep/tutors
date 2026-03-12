package audit

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	userCtx "github.com/it-chep/tutors.git/pkg/context"
)

type spyRepo struct {
	entries []Entry
	err     error
}

func (s *spyRepo) Create(_ context.Context, entry Entry) error {
	s.entries = append(s.entries, entry)
	return s.err
}

func TestMiddleware_LogsSuccessfulMappedRequest(t *testing.T) {
	repo := &spyRepo{}
	action := &Action{repo: repo}

	var handlerBody string

	router := chi.NewRouter()
	router.Use(action.Middleware())
	router.Post("/admin/students/{student_id}/archive", func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("read body in handler: %v", err)
		}
		handlerBody = string(body)
		w.WriteHeader(http.StatusNoContent)
	})

	req := httptest.NewRequest(http.MethodPost, "/admin/students/42/archive", strings.NewReader(`{"source":"ui"}`))
	req = req.WithContext(userCtx.WithUserID(req.Context(), 7))
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("unexpected status: got %d want %d", rec.Code, http.StatusNoContent)
	}
	if handlerBody != `{"source":"ui"}` {
		t.Fatalf("handler got body %q", handlerBody)
	}
	if len(repo.entries) != 1 {
		t.Fatalf("expected 1 audit entry, got %d", len(repo.entries))
	}

	entry := repo.entries[0]
	if entry.UserID != 7 {
		t.Fatalf("unexpected user id: got %d want 7", entry.UserID)
	}
	if entry.Description != "Архивация студента" {
		t.Fatalf("unexpected description: %q", entry.Description)
	}
	if entry.Action != "Архивация студента" {
		t.Fatalf("unexpected action: %q", entry.Action)
	}
	if entry.EntityName != "student" {
		t.Fatalf("unexpected entity name: %q", entry.EntityName)
	}
	if entry.EntityID == nil || *entry.EntityID != 42 {
		t.Fatalf("unexpected entity id: %#v", entry.EntityID)
	}
	if entry.Body == nil || *entry.Body != `{"source":"ui"}` {
		t.Fatalf("unexpected body: %#v", entry.Body)
	}
}

func TestMiddleware_DoesNotLogFailedRequest(t *testing.T) {
	repo := &spyRepo{}
	action := &Action{repo: repo}

	router := chi.NewRouter()
	router.Use(action.Middleware())
	router.Post("/admin/students/{student_id}/archive", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "bad request", http.StatusBadRequest)
	})

	req := httptest.NewRequest(http.MethodPost, "/admin/students/42/archive", strings.NewReader(`{"source":"ui"}`))
	req = req.WithContext(userCtx.WithUserID(req.Context(), 7))
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("unexpected status: got %d want %d", rec.Code, http.StatusBadRequest)
	}
	if len(repo.entries) != 0 {
		t.Fatalf("expected no audit entries, got %d", len(repo.entries))
	}
}

func TestMiddleware_SkipsUnmappedRouteAndKeepsBodyReadable(t *testing.T) {
	repo := &spyRepo{}
	action := &Action{repo: repo}

	var handlerBody string

	router := chi.NewRouter()
	router.Use(action.Middleware())
	router.Post("/admin/unmapped", func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("read body in handler: %v", err)
		}
		handlerBody = string(body)
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodPost, "/admin/unmapped", strings.NewReader(`{"hello":"world"}`))
	req = req.WithContext(userCtx.WithUserID(req.Context(), 7))
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("unexpected status: got %d want %d", rec.Code, http.StatusOK)
	}
	if handlerBody != `{"hello":"world"}` {
		t.Fatalf("handler got body %q", handlerBody)
	}
	if len(repo.entries) != 0 {
		t.Fatalf("expected no audit entries, got %d", len(repo.entries))
	}
}

func TestMiddleware_IgnoresAuditWriteErrors(t *testing.T) {
	repo := &spyRepo{err: errors.New("db unavailable")}
	action := &Action{repo: repo}

	router := chi.NewRouter()
	router.Use(action.Middleware())
	router.Post("/admin/students/{student_id}/archive", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})

	req := httptest.NewRequest(http.MethodPost, "/admin/students/42/archive", strings.NewReader(`{"source":"ui"}`))
	req = req.WithContext(userCtx.WithUserID(req.Context(), 7))
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("unexpected status: got %d want %d", rec.Code, http.StatusNoContent)
	}
	if len(repo.entries) != 1 {
		t.Fatalf("expected 1 attempted audit entry, got %d", len(repo.entries))
	}
}
