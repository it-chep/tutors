package create_comment

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/it-chep/tutors.git/internal/module/admin"
)

type Handler struct {
	adminModule *admin.Module
}

func NewHandler(adminModule *admin.Module) *Handler {
	return &Handler{
		adminModule: adminModule,
	}
}

func (h *Handler) Handle() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		studentIDStr := chi.URLParam(r, "student_id")
		studentID, err := strconv.ParseInt(studentIDStr, 10, 64)
		if err != nil {
			http.Error(w, "invalid student ID", http.StatusBadRequest)
			return
		}

		var req Request
		if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "failed to decode request: "+err.Error(), http.StatusInternalServerError)
			return
		}

		text := req.CleanText()
		if text == "" {
			http.Error(w, "comment text is required", http.StatusBadRequest)
			return
		}

		if err = h.adminModule.Actions.CreateComment.Do(ctx, studentID, text); err != nil {
			http.Error(w, "failed to create comment: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
