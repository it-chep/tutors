package delete_comment

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
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

		studentIDStr := chi.URLParam(r, "student_id")
		studentID, err := strconv.ParseInt(studentIDStr, 10, 64)
		if err != nil {
			http.Error(w, "invalid student ID", http.StatusBadRequest)
			return
		}

		commentIDStr := chi.URLParam(r, "comment_id")
		commentID, err := strconv.ParseInt(commentIDStr, 10, 64)
		if err != nil {
			http.Error(w, "invalid comment ID", http.StatusBadRequest)
			return
		}

		if err = h.adminModule.Actions.DeleteComment.Do(ctx, studentID, commentID); err != nil {
			http.Error(w, "failed to delete comment: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
