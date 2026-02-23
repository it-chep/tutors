package change_all_payment

import (
	"encoding/json"
	"net/http"

	"github.com/it-chep/tutors.git/internal/module/admin"
	"github.com/it-chep/tutors.git/internal/module/admin/dto"
	userCtx "github.com/it-chep/tutors.git/pkg/context"
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

		if dto.IsTutorRole(ctx) {
			http.Error(w, "Функционал недоступен", http.StatusForbidden)
			return
		}

		var req Request
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "failed to decode request: "+err.Error(), http.StatusBadRequest)
			return
		}

		adminID := userCtx.AdminIDFromContext(ctx)

		err := h.adminModule.Actions.ChangeAllStudentsPayment.Do(ctx, adminID, req.PaymentID)
		if err != nil {
			http.Error(w, "failed to change payment: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
