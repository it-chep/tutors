package update_tutor

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/it-chep/tutors.git/internal/module/admin"
	"github.com/it-chep/tutors.git/internal/module/admin/action/tutor/update_tutor/dto"
	indto "github.com/it-chep/tutors.git/internal/module/admin/dto"
)

type Handler struct {
	adminModule *admin.Module
}

func NewHandler(adminModule *admin.Module) *Handler {
	return &Handler{
		adminModule: adminModule,
	}
}

type Request struct {
	FullName        string `json:"full_name"`
	Phone           string `json:"phone"`
	Tg              string `json:"tg"`
	CostPerHour     string `json:"cost_per_hour"`
	SubjectID       int64  `json:"subject_id"`
	TgAdminUsername string `json:"tg_admin_username"`
}

func (req Request) ToDto() dto.UpdateRequest {
	return dto.UpdateRequest{
		FullName:        req.FullName,
		Phone:           req.Phone,
		Tg:              req.Tg,
		CostPerHour:     req.CostPerHour,
		SubjectID:       req.SubjectID,
		TgAdminUsername: req.TgAdminUsername,
	}
}

func (h *Handler) Handle() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		if indto.IsTutorRole(ctx) {
			http.Error(w, "authorization required", http.StatusForbidden)
			return
		}

		var req Request
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "failed to decode request: "+err.Error(), http.StatusBadRequest)
			return
		}

		tutorIDStr := chi.URLParam(r, "tutor_id")
		tutorID, err := strconv.ParseInt(tutorIDStr, 10, 64)
		if err != nil {
			http.Error(w, "invalid tutor ID", http.StatusBadRequest)
			return
		}

		err = h.adminModule.Actions.UpdateTutor.Do(ctx, tutorID, req.ToDto())
		if err != nil {
			http.Error(w, "failed to update tutor: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
