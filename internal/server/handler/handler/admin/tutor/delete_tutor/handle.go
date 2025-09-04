package delete_tutor

import (
	"github.com/go-chi/chi/v5"
	"github.com/it-chep/tutors.git/internal/module/admin"
	"net/http"
	"strconv"
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

		tutorIDStr := chi.URLParam(r, "tutor_id")
		tutorID, err := strconv.ParseInt(tutorIDStr, 10, 64)
		if err != nil {
			http.Error(w, "invalid tutor ID", http.StatusBadRequest)
			return
		}

		err = h.adminModule.Actions.DeleteTutor.Do(ctx, tutorID)
		if err != nil {
			http.Error(w, "failed to delete tutor data: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
