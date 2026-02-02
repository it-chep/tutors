package get_admin_payments

import (
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

func (h *Handler) prepareResponse(payments []dto.Payment) Response {
	return Response{
		Payments: lo.Map(payments, func(item dto.Payment, index int) Payment {
			return Payment{
				ID:   item.ID,
				Name: item.String(),
			}
		}),
	}
}
