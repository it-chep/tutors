package get_all_finance

import (
	"context"
	"encoding/json"
	"github.com/it-chep/tutors.git/internal/module/admin/action/get_all_finance/dto"
	"net/http"

	indto "github.com/it-chep/tutors.git/internal/module/admin/dto"
	userCtx "github.com/it-chep/tutors.git/pkg/context"

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
			http.Error(w, "у репетитора нет на это прав ", http.StatusUnauthorized)
			return
		}

		// суперадмин отправит ID админа в теле
		adminID := req.AdminID
		if indto.IsAdminRole(ctx) || indto.IsAssistantRole(ctx) {
			adminID = userCtx.AdminIDFromContext(ctx)
		}

		finance, err := h.adminModule.Actions.GetAllFinance.Do(ctx, req.From, req.To, adminID)
		if err != nil {
			http.Error(w, "failed to get user data: "+err.Error(), http.StatusInternalServerError)
			return
		}

		response := h.prepareResponse(ctx, finance)

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, "failed to encode response: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func (h *Handler) prepareResponse(ctx context.Context, finance dto.GetAllFinanceDto) Response {
	if indto.IsAssistantRole(ctx) {
		return Response{
			Finance: Finance{
				Salary: finance.TutorsInfo.Salary,
				Hours:  finance.TutorsInfo.Hours,
			},
		}
	}

	return Response{
		Finance: Finance{
			Profit:   finance.Profit,
			CashFlow: finance.CashFlow,
			Debt:     finance.Debt,
			Salary:   finance.TutorsInfo.Salary,
			Hours:    finance.TutorsInfo.Hours,
		},
	}
}
