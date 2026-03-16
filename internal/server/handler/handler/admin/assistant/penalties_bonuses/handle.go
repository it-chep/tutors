package penalties_bonuses

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/it-chep/tutors.git/internal/module/admin"
	"github.com/it-chep/tutors.git/internal/module/admin/action/accrual"
	"github.com/it-chep/tutors.git/internal/module/admin/dto"
	"github.com/samber/lo"
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

		assistantID, err := strconv.ParseInt(chi.URLParam(r, "assistant_id"), 10, 64)
		if err != nil {
			http.Error(w, "invalid assistant ID", http.StatusBadRequest)
			return
		}

		var req Request
		if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "failed to decode request: "+err.Error(), http.StatusInternalServerError)
			return
		}

		if req.IsCreate() {
			actualType, err := req.ToActualType()
			if err != nil {
				http.Error(w, "invalid type", http.StatusBadRequest)
				return
			}

			err = h.adminModule.Actions.Accruals.Create(ctx, accrual.CreateRequest{
				TargetUserID: assistantID,
				TargetRoleID: dto.AssistantRole,
				Type:         actualType,
				Amount:       req.Amount,
				Comment:      req.Comment,
			})
			if err != nil {
				http.Error(w, "failed to create accrual: "+err.Error(), http.StatusInternalServerError)
				return
			}
			return
		}

		from, to, err := req.ToTime()
		if err != nil {
			http.Error(w, "invalid time", http.StatusBadRequest)
			return
		}

		accrualsList, summary, err := h.adminModule.Actions.Accruals.List(ctx, assistantID, dto.AssistantRole, from, to)
		if err != nil {
			http.Error(w, "failed to get penalties and bonuses: "+err.Error(), http.StatusInternalServerError)
			return
		}

		filtered := lo.Filter(accrualsList, func(item dto.Accrual, _ int) bool {
			return item.ActualTypeID == dto.AccrualActualTypePenalty || item.ActualTypeID == dto.AccrualActualTypeBonus
		})

		response := Response{
			Items: lo.Map(filtered, func(item dto.Accrual, _ int) Item {
				return Item{
					ID:       item.ID,
					Type:     item.ActualTypeID.String(),
					Amount:   item.Amount.String(),
					Comment:  item.Comment,
					ActualAt: item.ActualAt.Format(time.DateTime),
					IsPaid:   item.IsPaid,
				}
			}),
			Summary: Summary{
				Penalties: summary.Penalties.String(),
				Bonuses:   summary.Bonuses.String(),
				Payable:   summary.Payable.String(),
				Unpaid:    summary.Unpaid.String(),
			},
		}

		w.Header().Set("Content-Type", "application/json")
		if err = json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, "failed to encode response: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
