package get_tg_admins_usernames

import (
	"encoding/json"
	"github.com/it-chep/tutors.git/internal/module/admin"
	"github.com/it-chep/tutors.git/internal/module/admin/dto"
	userCtx "github.com/it-chep/tutors.git/pkg/context"
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

		if dto.IsTutorRole(ctx) {
			http.Error(w, "Нет прав на это", http.StatusForbidden)
			return
		}

		adminID := userCtx.AdminIDFromContext(ctx)

		usernames, err := h.adminModule.Actions.GetTgAdminsUsernames.Do(ctx, adminID)
		response := Response{
			Usernames: lo.Ternary(len(usernames) > 0, usernames, []string{}),
		}

		w.Header().Set("Content-Type", "application/json")
		if err = json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, "failed to encode response: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
