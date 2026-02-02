package change_student_payment

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/it-chep/tutors.git/internal/module/admin"
	"github.com/it-chep/tutors.git/internal/module/admin/dto"
	userCtx "github.com/it-chep/tutors.git/pkg/context"
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

		studentIDStr := chi.URLParam(r, "student_id")
		studentID, err := strconv.ParseInt(studentIDStr, 10, 64)
		if err != nil {
			http.Error(w, "invalid user ID", http.StatusBadRequest)
			return
		}

		var req Request
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "failed to decode request: "+err.Error(), http.StatusInternalServerError)
			return
		}

		if dto.IsTutorRole(ctx) {
			http.Error(w, "", http.StatusForbidden)
			return
		}

		adminID := userCtx.AdminIDFromContext(ctx)

		err = h.adminModule.Actions.ChangeStudentPayment.Do(ctx, adminID, studentID, req.NewPaymentID)
		if err != nil {
			http.Error(w, "failed to create student data: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
