package get_assistant_by_id

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/it-chep/tutors.git/internal/module/admin/dto"
	"github.com/samber/lo"
	"net/http"
	"strconv"
	"time"

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

		assistantIDStr := chi.URLParam(r, "assistant_id")
		assistantID, err := strconv.ParseInt(assistantIDStr, 10, 64)
		if err != nil {
			http.Error(w, "invalid assistant ID", http.StatusBadRequest)
			return
		}

		baseData, err := h.adminModule.Actions.GetAdminByID.Do(ctx, assistantID)
		if err != nil {
			http.Error(w, "failed to get assistant data: "+err.Error(), http.StatusInternalServerError)
			return
		}
		tgs, err := h.adminModule.Actions.GetAssistantAvailableTGs.Do(ctx, assistantID)
		if err != nil {
			http.Error(w, "failed to get assistant data: "+err.Error(), http.StatusInternalServerError)
			return
		}
		response := h.prepareResponse(baseData, tgs)

		w.Header().Set("Content-Type", "application/json")
		if err = json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, "failed to encode response: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func (h *Handler) prepareResponse(admin dto.User, tgs dto.TgAdminUsernames) Response {
	return Response{
		Assistant: Assistant{
			ID:       admin.ID,
			FullName: admin.FullName,
			Tg:       admin.Tg,
			Phone:    admin.Phone,
			AvailableTgs: lo.Map(tgs, func(item dto.TgAdminUsername, _ int) TgAdminUsername {
				return TgAdminUsername{
					ID:   item.ID,
					Name: item.Name,
				}
			}),
			CreatedAt: admin.CreatedAt.Format(time.DateTime),
		},
	}
}
