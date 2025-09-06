package create_student

import (
	"encoding/json"
	"net/http"

	"github.com/it-chep/tutors.git/internal/module/admin"
	"github.com/it-chep/tutors.git/internal/module/admin/action/student/create_student/dto"
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

		var req Request
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "failed to decode request: "+err.Error(), http.StatusInternalServerError)
			return
		}

		err := h.adminModule.Actions.CreateStudent.Do(ctx, dto.CreateRequest{
			FirstName:  req.FirstName,
			LastName:   req.LastName,
			MiddleName: req.MiddleName,

			Phone:       req.Phone,
			Tg:          req.Tg,
			CostPerHour: req.CostPerHour,
			SubjectID:   req.SubjectID,
			TutorID:     req.TutorID,

			ParentFullName: req.ParentFullName,
			ParentPhone:    req.ParentPhone,
			ParentTg:       req.ParentTg,
		})
		if err != nil {
			http.Error(w, "failed to create student data: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
