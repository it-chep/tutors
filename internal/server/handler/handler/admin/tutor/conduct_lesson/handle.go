package conduct_lesson

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/it-chep/tutors.git/internal/module/admin/dto"
	userCtx "github.com/it-chep/tutors.git/pkg/context"

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
		var (
			req     Request
			tutorID int64
			ctx     = r.Context()
		)

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "failed to decode request: "+err.Error(), http.StatusInternalServerError)
			return
		}

		if dto.IsTutorRole(ctx) {
			tutorID = userCtx.GetTutorID(ctx)
		}

		createdTime, err := time.Parse(time.DateTime, req.Date)
		if err != nil {
			http.Error(w, "failed to parse date: "+err.Error(), http.StatusBadRequest)
		}

		err = h.adminModule.Actions.ConductLesson.Do(ctx, tutorID, req.StudentID, req.Duration, createdTime)
		if err != nil {
			http.Error(w, "failed to conduct lesson: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
