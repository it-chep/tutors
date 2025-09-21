package get_all_finance

import (
	"encoding/json"
	indto "github.com/it-chep/tutors.git/internal/module/admin/dto"
	"net/http"

	"github.com/it-chep/tutors.git/internal/module/admin/action/get_all_finance/dto"

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

		var req Request
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "failed to decode request: "+err.Error(), http.StatusInternalServerError)
			return
		}

		if indto.IsTutorRole(ctx) {
			http.Error(w, "authorization required", http.StatusUnauthorized)
			return
		}

		finance, err := h.adminModule.Actions.GetAllFinance.Do(ctx, req.From, req.To)
		if err != nil {
			http.Error(w, "failed to get user data: "+err.Error(), http.StatusInternalServerError)
			return
		}

		response := h.prepareResponse(finance)

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, "failed to encode response: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func (h *Handler) prepareResponse(finance dto.GetAllFinanceDto) Response {
	return Response{
		Finance: Finance{
			Profit:            finance.Profit,
			CashFlow:          finance.CashFlow,
			Conversion:        finance.Conversion,
			LessonsCount:      finance.CountLessons,
			CountBaseLessons:  finance.CountBaseLessons,
			CountTrialLessons: finance.CountTrialLessons,
		},
	}
}
