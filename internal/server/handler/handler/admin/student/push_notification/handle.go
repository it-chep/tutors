package push_notification

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

		studentIDStr := chi.URLParam(r, "student_id")
		studentID, err := strconv.ParseInt(studentIDStr, 10, 64)
		if err != nil {
			http.Error(w, "invalid user ID", http.StatusBadRequest)
			return
		}

		err = h.adminModule.Actions.PushNotification.Do(ctx, studentID)
		if err != nil {
			http.Error(w, "failed to push notification: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
