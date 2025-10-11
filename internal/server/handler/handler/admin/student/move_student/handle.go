package move_student

import (
	"encoding/json"
	"github.com/it-chep/tutors.git/internal/module/admin"
	"github.com/it-chep/tutors.git/internal/module/admin/dto"
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

		if dto.IsTutorRole(ctx) {
			http.Error(w, "", http.StatusForbidden)
			return
		}

		err := h.adminModule.Actions.MoveStudents.Do(ctx, req.OldTutorID, req.NewTutorID, req.StudentID)
		if err != nil {
			http.Error(w, "failed to create student data: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
