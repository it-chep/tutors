package push_all_debitors

import (
	"github.com/it-chep/tutors.git/internal/module/admin"
	userCtx "github.com/it-chep/tutors.git/pkg/context"
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

		err := h.adminModule.Actions.PushAllDebitors.Do(ctx, userCtx.AdminIDFromContext(ctx))
		if err != nil {
			http.Error(w, "failed to get user data: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
