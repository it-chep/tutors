package get_all_finance

import (
	"github.com/it-chep/tutors.git/internal/module/admin"
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
		_ = r.Context()

	}
}

func (h *Handler) prepareResponse() Response {
	return Response{}
}
