package get_tutor_finance

import (
	"encoding/json"
	"github.com/it-chep/tutors.git/internal/module/admin"
	"github.com/it-chep/tutors.git/internal/module/admin/action/tutor/create_tutor/dto"
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

		if r.Method != http.MethodPost {
			http.Error(w, "", http.StatusMethodNotAllowed)
			return
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "failed to decode request: "+err.Error(), http.StatusInternalServerError)
			return
		}

		err := h.adminModule.Actions.CreateInformationPost.Do(ctx, dto.Request{
			PostName:      req.PostName,
			ThemeID:       req.ThemeID,
			Order:         req.Order,
			MediaID:       req.MediaID,
			ContentTypeID: req.ContentTypeID,
			Message:       req.Message,
		})
		if err != nil {
			http.Error(w, "failed to create information post: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
