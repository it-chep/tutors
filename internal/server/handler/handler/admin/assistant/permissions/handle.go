package permissions

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/it-chep/tutors.git/internal/module/admin"
	permissionsAction "github.com/it-chep/tutors.git/internal/module/admin/action/assistant/permissions"
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

		assistantID, err := strconv.ParseInt(chi.URLParam(r, "assistant_id"), 10, 64)
		if err != nil {
			http.Error(w, "invalid assistant ID", http.StatusBadRequest)
			return
		}

		var req Request
		if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "failed to decode request: "+err.Error(), http.StatusInternalServerError)
			return
		}

		err = h.adminModule.Actions.AssistantPermissions.Update(ctx, assistantID, permissionsAction.UpdateRequest{
			CanViewContracts:      req.CanViewContracts,
			CanPenalizeAssistants: req.CanPenalizeAssistants,
		})
		if err != nil {
			http.Error(w, "failed to update permissions: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
