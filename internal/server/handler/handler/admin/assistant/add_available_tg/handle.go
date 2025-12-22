package add_available_tg

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/it-chep/tutors.git/internal/module/admin"
	"net/http"
	"strconv"
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

		var req Request
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "failed to decode request: "+err.Error(), http.StatusInternalServerError)
			return
		}

		err = h.adminModule.Actions.AddAvailableTg.Do(ctx, assistantID, req.AvailableTg)
		if err != nil {
			http.Error(w, "failed to get assistant data: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
