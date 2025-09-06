package get_all_subjects

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

		subjects, err := h.adminModule.Actions.GetAllSubjects.Do(ctx)
		if err != nil {
			http.Error(w, "failed to get subjects data: "+err.Error(), http.StatusInternalServerError)
			return
		}
		response := h.prepareResponse(subjects)

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, "failed to encode response: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func (h *Handler) prepareResponse(subjects []dto.Subject) Response {
	return Response{
		Subjects: lo.Map(subjects, func(item dto.Subject, index int) Subject {
			return Subject{
				ID:   item.ID,
				Name: item.Name,
			}
		}),
	}
}
