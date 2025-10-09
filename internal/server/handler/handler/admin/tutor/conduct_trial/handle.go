package conduct_trial

import (
	"encoding/json"
	"github.com/it-chep/tutors.git/internal/module/admin/dto"
	userCtx "github.com/it-chep/tutors.git/pkg/context"
	"net/http"

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

		var (
			req     Request
			tutorID int64
		)
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "failed to decode request: "+err.Error(), http.StatusInternalServerError)
			return
		}

		if dto.IsTutorRole(ctx) {
			tutorID = userCtx.GetTutorID(ctx)
		}

		err := h.adminModule.Actions.ConductTrial.Do(ctx, tutorID, req.StudentID)
		if err != nil {
			http.Error(w, "failed to conduct lesson: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
