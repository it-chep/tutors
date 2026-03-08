package get_comments

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

		studentIDStr := chi.URLParam(r, "student_id")
		studentID, err := strconv.ParseInt(studentIDStr, 10, 64)
		if err != nil {
			http.Error(w, "invalid student ID", http.StatusBadRequest)
			return
		}

		comments, err := h.adminModule.Actions.GetComments.Do(ctx, studentID)
		if err != nil {
			http.Error(w, "failed to get comments: "+err.Error(), http.StatusInternalServerError)
			return
		}

		response := h.prepareResponse(comments)

		w.Header().Set("Content-Type", "application/json")
		if err = json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, "failed to encode response: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func (h *Handler) prepareResponse(comments []dto.Comment) Response {
	return Response{
		Comments: lo.Map(comments, func(item dto.Comment, _ int) Comment {
			return Comment{
				ID:             item.ID,
				Text:           item.Text,
				AuthorFullName: item.AuthorFullName,
				CreatedAt:      item.CreatedAt.Format(time.DateTime),
			}
		}),
		CommentsCount: int64(len(comments)),
	}
}
