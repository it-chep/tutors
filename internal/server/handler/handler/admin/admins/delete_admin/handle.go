package delete_admin

import (
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"

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

		err = h.adminModule.Actions.DeleteAdmin.Do(ctx, adminID)
		if err != nil {
			http.Error(w, "failed to create student data: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
