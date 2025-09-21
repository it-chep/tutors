package get_admin_by_id

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/it-chep/tutors.git/internal/module/admin/dto"

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

		adminIDStr := chi.URLParam(r, "admin_id")
		adminID, err := strconv.ParseInt(adminIDStr, 10, 64)
		if err != nil {
			http.Error(w, "invalid admin ID", http.StatusBadRequest)
			return
		}

		baseData, err := h.adminModule.Actions.GetAdminByID.Do(ctx, adminID)
		if err != nil {
			http.Error(w, "failed to get admins data: "+err.Error(), http.StatusInternalServerError)
			return
		}

		response := h.prepareResponse(baseData)

		w.Header().Set("Content-Type", "application/json")
		if err = json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, "failed to encode response: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func (h *Handler) prepareResponse(admin dto.User) Response {
	return Response{
		Admin: Admin{
			ID:       admin.ID,
			FullName: admin.FullName,
			Tg:       admin.Tg,
		},
	}
}
