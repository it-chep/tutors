package delete_assistant

import (
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

		assistantIDStr := chi.URLParam(r, "assistant_id")
		assistantID, err := strconv.ParseInt(assistantIDStr, 10, 64)
		if err != nil {
			http.Error(w, "invalid assistant ID", http.StatusBadRequest)
			return
		}

		err = h.adminModule.Actions.DeleteAdmin.Do(ctx, assistantID, dto.AssistantRole)
		if err != nil {
			http.Error(w, "failed to create student data: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
