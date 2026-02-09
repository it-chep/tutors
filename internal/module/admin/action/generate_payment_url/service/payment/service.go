package payment

import (
	"context"
	"github.com/it-chep/tutors.git/internal/config"
	"github.com/it-chep/tutors.git/internal/pkg/logger"
)

type Repository interface {
	PaymentAndAdminByStudent(ctx context.Context, studentID int64) (int64, int64, error)
}

type Service struct {
	dal            Repository
	paymentByAdmin config.PaymentsByAdmin
}

func NewService(dal Repository, paymentByAdmin config.PaymentsByAdmin) *Service {
	return &Service{
		dal:            dal,
		paymentByAdmin: paymentByAdmin,
	}
}

// GetPayment получаем платежку по ID админа и платежки
func (s *Service) GetPayment(ctx context.Context, studentID int64) config.AdminPayment {
	adminID, paymentID, err := s.dal.PaymentAndAdminByStudent(ctx, studentID)
	if err != nil {
		logger.Error(ctx, "ошибка при получении админа от тутора родителя", err)
		return config.AdminPayment{}
	}

	return s.paymentByAdmin.Payment(adminID, paymentID)
}
