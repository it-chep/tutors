package get_tutor_finance

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/it-chep/tutors.git/internal/module/admin"
	"github.com/it-chep/tutors.git/internal/module/admin/dto"
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

		tutorIDStr := chi.URLParam(r, "tutor_id")
		tutorID, err := strconv.ParseInt(tutorIDStr, 10, 64)
		if err != nil {
			http.Error(w, "invalid tutor_id", http.StatusBadRequest)
			return
		}

		if dto.IsTutorRole(ctx) {
			http.Error(w, "authorization required", http.StatusUnauthorized)
			return
		}

		var req Request
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "failed to decode request: "+err.Error(), http.StatusInternalServerError)
			return
		}

		from, to, err := req.ToTime()
		if err != nil {
			http.Error(w, "invalid time", http.StatusBadRequest)
			return
		}

		baseData, err := h.adminModule.Actions.GetTutorFinance.Do(ctx, tutorID, from, to)
		if err != nil {
			http.Error(w, "failed to get finance data: "+err.Error(), http.StatusInternalServerError)
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

func (h *Handler) prepareResponse(financeInfo dto.TutorFinance) Response {
	return Response{
		Finance: Finance{
			HoursCount: financeInfo.HoursCount,
			Wages:      financeInfo.Wages.String(),
			Amount:     financeInfo.Amount.String(),
		},
	}
}
