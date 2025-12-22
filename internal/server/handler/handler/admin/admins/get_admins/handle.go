package get_admins

import (
	"encoding/json"
	"net/http"

	"github.com/it-chep/tutors.git/internal/module/admin/dto"
	"github.com/samber/lo"

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

		baseData, err := h.adminModule.Actions.GetAdmins.Do(ctx, dto.AdminRole)
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

func (h *Handler) prepareResponse(admins []dto.User) Response {
	return Response{
		Admins: lo.Map(admins, func(item dto.User, index int) Admin {
			return Admin{
				ID:       item.ID,
				FullName: item.FullName,
				Tg:       item.Tg,
			}
		}),
	}
}
