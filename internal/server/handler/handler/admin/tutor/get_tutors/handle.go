package get_tutors

import (
	"encoding/json"
	"github.com/it-chep/tutors.git/internal/module/admin"
	"github.com/it-chep/tutors.git/internal/module/admin/dto"
	"github.com/samber/lo"
	"net/http"
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

		//tutorIDStr := r.URL.Query().Get("admin_id")
		//_, err := strconv.ParseInt(tutorIDStr, 10, 64)
		//if err != nil {
		//	http.Error(w, "invalid admin ID", http.StatusBadRequest)
		//	return
		//}

		baseData, err := h.adminModule.Actions.GetTutors.Do(ctx)
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

func (h *Handler) prepareResponse(students []dto.Tutor) Response {
	return Response{
		Tutors: lo.Map(students, func(item dto.Tutor, index int) Tutor {
			return Tutor{
				ID:                 item.ID,
				FullName:           item.FullName,
				Tg:                 item.Tg,
				HasBalanceNegative: false,
				HasOnlyTrial:       false,
				HasNewBie:          false,
			}
		}),
	}
}
