package tutor_by_id

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/it-chep/tutors.git/internal/module/admin"
	"github.com/it-chep/tutors.git/internal/module/admin/dto"
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

		baseData, err := h.adminModule.Actions.GetTutorByID.Do(ctx, tutorID)
		if err != nil {
			http.Error(w, "failed to get tutor data: "+err.Error(), http.StatusInternalServerError)
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

func (h *Handler) prepareResponse(tutor dto.Tutor) Response {
	return Response{
		Tutor: Tutor{
			ID:          tutor.ID,
			FullName:    tutor.FullName,
			Phone:       tutor.Phone,
			Tg:          tutor.Tg,
			CostPerHour: tutor.CostPerHour,
			SubjectName: tutor.SubjectName,
		},
	}
}
