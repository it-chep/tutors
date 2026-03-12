package accruals

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/it-chep/tutors.git/internal/module/admin"
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

		from, to, err := req.ToTime()
		if err != nil {
			http.Error(w, "invalid time", http.StatusBadRequest)
			return
		}

		accruals, summary, err := h.adminModule.Actions.Accruals.List(ctx, tutorID, dto.TutorRole, from, to)
		if err != nil {
			http.Error(w, "failed to get accruals: "+err.Error(), http.StatusInternalServerError)
			return
		}

		response := Response{
			Accruals: lo.Map(accruals, func(item dto.Accrual, _ int) Accrual {
				return Accrual{
					ID:         item.ID,
					ActualType: item.ActualTypeID.String(),
					Amount:     item.Amount.String(),
					Comment:    item.Comment,
					LessonID:   item.LessonID,
					ActualAt:   item.ActualAt.Format(time.DateTime),
					IsPaid:     item.IsPaid,
				}
			}),
			Summary: Summary{
				Lessons:   summary.Lessons.String(),
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
