package get_all_finance_by_tgs

import (
	"encoding/json"
	indto "github.com/it-chep/tutors.git/internal/module/admin/dto"
	"github.com/it-chep/tutors.git/internal/pkg/convert"
	userCtx "github.com/it-chep/tutors.git/pkg/context"
	"net/http"

	"github.com/it-chep/tutors.git/internal/module/admin/action/get_all_finance_by_tgs/dto"

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

		from, to, err := convert.StringsIntervalToTime(req.From, req.To)
		if err != nil {
			http.Error(w, "failed to convert from time: "+err.Error(), http.StatusBadRequest)
			return
		}
		// суперадмин отправит ID админа в теле
		adminID := req.AdminID
		if indto.IsAdminRole(ctx) {
			adminID = userCtx.UserIDFromContext(ctx)
		}

		finance, err := h.adminModule.Actions.GetAllFinanceByTGs.Do(ctx, dto.Request{
			AdminID:     adminID,
			TgUsernames: req.TgUsernames,
			From:        from,
			To:          to,
		})
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
			Profit:   finance.Profit,
			CashFlow: finance.CashFlow,
			Debt:     finance.Debt,
			Salary:   finance.TutorsInfo.Salary,
			Hours:    finance.TutorsInfo.Hours,
		},
	}
}
