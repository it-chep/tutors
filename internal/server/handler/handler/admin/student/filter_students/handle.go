package filter_students

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

		// todo надо ли такое делать в репетиторе ?
		//tutorIDStr := r.URL.Query().Get("tutor_id")
		//tutorID, err := strconv.ParseInt(tutorIDStr, 10, 64)
		//if err != nil {
		//	tutorID = 0
		//}
		//if dto.IsTutorRole(ctx) {
		//	tutorID = userCtx.GetTutorID(ctx)
		//}

		var req Request
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "failed to decode request: "+err.Error(), http.StatusInternalServerError)
			return
		}

		students, err := h.adminModule.Actions.FilterStudents.Do(ctx, req.ToFilterRequest())
		if err != nil {
			http.Error(w, "failed to get user data: "+err.Error(), http.StatusInternalServerError)
			return
		}

		response := h.prepareResponse(students)

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
				Balance:             item.Balance.String(),
			}
		}),
		StudentsCount: int64(len(students)),
	}
}
