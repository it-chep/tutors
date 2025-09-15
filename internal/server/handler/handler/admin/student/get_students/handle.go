package get_students

import (
	"encoding/json"
	"net/http"
	"strconv"

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

		tutorIDStr := r.URL.Query().Get("tutor_id")
		tutorID, err := strconv.ParseInt(tutorIDStr, 10, 64)
		if err != nil {
			tutorID = 0
		}

		baseData, err := h.adminModule.Actions.GetStudents.Do(ctx, tutorID)
		if err != nil {
			http.Error(w, "failed to get user data: "+err.Error(), http.StatusInternalServerError)
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

func (h *Handler) prepareResponse(students []dto.Student) Response {
	return Response{
		Students: lo.Map(students, func(item dto.Student, index int) Student {
			return Student{
				ID:                  item.ID,
				FirstName:           item.FirstName,
				LastName:            item.LastName,
				MiddleName:          item.MiddleName,
				ParentFullName:      item.ParentFullName,
				Tg:                  item.Tg,
				IsOnlyTrialFinished: item.IsOnlyTrialFinished,
				IsBalanceNegative:   item.IsBalanceNegative,
				IsNewbie:            item.IsNewbie,
			}
		}),
	}
}
