package get_tutors

import (
	"encoding/json"
	"net/http"
	"strconv"

	userCtx "github.com/it-chep/tutors.git/pkg/context"

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

		adminIDStr := r.URL.Query().Get("admin_id")
		adminID, err := strconv.ParseInt(adminIDStr, 10, 64)
		if err != nil {
			adminID = 0
		}

		if adminID == 0 && dto.IsAdminRole(ctx) {
			adminID = userCtx.UserIDFromContext(ctx)
		}

		baseData, err := h.adminModule.Actions.GetTutors.Do(ctx, adminID)
		if err != nil {
			http.Error(w, "failed to get tutors data: "+err.Error(), http.StatusInternalServerError)
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

func (h *Handler) prepareResponse(tutors []dto.Tutor) Response {
	return Response{
		Tutors: lo.Map(tutors, func(item dto.Tutor, index int) Tutor {
			return Tutor{
				ID:                 item.ID,
				FullName:           item.FullName,
				Tg:                 item.Tg,
				HasBalanceNegative: item.HasBalanceNegative,
				HasOnlyTrial:       item.HasOnlyTrial,
				HasNewBie:          item.HasNewBie,
			}
		}),
		TutorsCount: int64(len(tutors)),
	}
}
