package create_assistant

import (
	"encoding/json"
	"net/http"

	"github.com/it-chep/tutors.git/internal/module/admin/action/admin/create_admin/dto"
	dto2 "github.com/it-chep/tutors.git/internal/module/admin/dto"

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

		var req Request
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "failed to decode request: "+err.Error(), http.StatusInternalServerError)
			return
		}

		err := h.adminModule.Actions.CreateAdmin.Do(ctx, dto.CreateRequest{
			FullName: req.FullName,
			Tg:       req.Tg,
			Phone:    req.Phone,
			Email:    req.Email,
			Role:     dto2.AssistantRole,
		})
		if err != nil {
			http.Error(w, "failed to create student data: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
