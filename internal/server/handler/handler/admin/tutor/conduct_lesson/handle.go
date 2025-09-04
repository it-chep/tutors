package conduct_lesson

import (
	"encoding/json"
	"github.com/it-chep/tutors.git/internal/module/admin"
	"net/http"
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

		var req Request
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "failed to decode request: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// todo tutor_ID
		err := h.adminModule.Actions.ConductLesson.Do(ctx, 0, req.StudentID, req.Duration)
		if err != nil {
			http.Error(w, "failed to conduct lesson: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
