package create_tutor

import (
	"encoding/json"
	"net/http"

	userCtx "github.com/it-chep/tutors.git/pkg/context"

	"github.com/it-chep/tutors.git/internal/module/admin"
	"github.com/it-chep/tutors.git/internal/module/admin/action/tutor/create_tutor/dto"
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

		if r.Method != http.MethodPost {
			http.Error(w, "", http.StatusMethodNotAllowed)
			return
		}

		var req Request
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "failed to decode request: "+err.Error(), http.StatusInternalServerError)
			return
		}

		adminID := userCtx.UserIDFromContext(ctx)

		err := h.adminModule.Actions.CreateTutor.Do(ctx, dto.Request{
			FullName:    req.FullName,
			Phone:       req.Phone,
			Tg:          req.Tg,
			CostPerHour: req.CostPerHour,
			SubjectID:   req.SubjectID,
		}, adminID)
		if err != nil {
			http.Error(w, "failed to create tutor: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
