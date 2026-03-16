package payouts

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/it-chep/tutors.git/internal/module/admin"
	"github.com/it-chep/tutors.git/internal/module/admin/action/payout"
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

		tutorID, err := strconv.ParseInt(chi.URLParam(r, "tutor_id"), 10, 64)
		if err != nil {
			http.Error(w, "invalid tutor ID", http.StatusBadRequest)
			return
		}

		var req Request
		if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "failed to decode request: "+err.Error(), http.StatusInternalServerError)
			return
		}

		payoutID, err := h.adminModule.Actions.Payouts.Create(ctx, payout.CreateRequest{
			TutorID: tutorID,
			Amount:  req.Amount,
			Comment: req.Comment,
		})
		if err != nil {
			http.Error(w, "failed to create payout: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]string{
			"id": payoutID.String(),
		})
	}
}
