package validator

import (
	"context"
	"github.com/it-chep/tutors.git/internal/module/admin/dto"
	"github.com/pkg/errors"
)

const (
	maxCountTransactionInMinute = 3
	maxTransactionAmount        = 30_000
)

type Repository interface {
	GetStudentByIDAndPaymentUUID(ctx context.Context, studentID int64, studentUUID string) (dto.Student, error)
	CountLastMinuteTransactions(ctx context.Context, studentID int64) (int64, error)
}

type Service struct {
	dal Repository
}

func New(dal Repository) *Service {
	return &Service{
		dal: dal,
	}
}

func (s *Service) Validate(ctx context.Context, studentID, amount int64, studentUUID string) error {
	_, err := s.dal.GetStudentByIDAndPaymentUUID(ctx, studentID, studentUUID)
	if err != nil {
		return err
	}

	transactionsCount, err := s.dal.CountLastMinuteTransactions(ctx, studentID)
	if err != nil {
		return err
	}

	if transactionsCount > maxCountTransactionInMinute {
		return errors.New("Превышен лимит по количеству транзакций")
	}

	if amount > maxTransactionAmount {
		return errors.New("Превышен лимит по вводимой сумме")
	}

	return nil

}
