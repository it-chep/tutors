package get_admin_payments

import (
	"encoding/json"
	userCtx "github.com/it-chep/tutors.git/pkg/context"
	"net/http"
)

type Payment struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}
type Response struct {
	Payments []Payment `json:"payments"`
}

func (h *Handler) Handle() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		adminID := userCtx.AdminIDFromContext(ctx)

		payments, err := h.adminModule.Actions.GetAdminAvailablePayments.Do(ctx, adminID)
		if err != nil {
			http.Error(w, "failed to get admin payments: "+err.Error(), http.StatusInternalServerError)
			return
		}

		response := h.prepareResponse(payments)

		w.Header().Set("Content-Type", "application/json")
		if err = json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, "failed to encode response: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
