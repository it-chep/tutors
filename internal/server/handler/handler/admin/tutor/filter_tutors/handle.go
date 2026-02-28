package filter_tutors

import (
	"encoding/json"
	"net/http"

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

		var req Request
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "failed to decode request: "+err.Error(), http.StatusInternalServerError)
			return
		}

		filterReq := req.ToFilterRequest()
		tutors, err := h.adminModule.Actions.FilterTutors.Do(ctx, filterReq)
		if err != nil {
			http.Error(w, "failed to get tutors data: "+err.Error(), http.StatusInternalServerError)
			return
		}

		response := h.prepareResponse(tutors)

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
				TgAdminUsername:    item.TgAdminUsername,
			}
		}),
		TutorsCount: int64(len(tutors)),
	}
}
