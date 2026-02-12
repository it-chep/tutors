package payments_generator

import (
	"context"
	"fmt"
	"github.com/it-chep/tutors.git/internal/config"
	dtoInternal "github.com/it-chep/tutors.git/internal/dto"
	"github.com/it-chep/tutors.git/internal/module/admin/action/generate_payment_url/dto"
	alfadto "github.com/it-chep/tutors.git/internal/pkg/alpha/dto"
	"github.com/it-chep/tutors.git/internal/pkg/logger"
	tbankDto "github.com/it-chep/tutors.git/internal/pkg/tbank/dto"
	tochkaDto "github.com/it-chep/tutors.git/internal/pkg/tochka/dto"
	"github.com/pkg/errors"
)

type Repository interface {
	DropTransaction(ctx context.Context, internalTransactionID string) error
	PhoneByStudent(ctx context.Context, studentID int64) (string, error)
}

type Service struct {
	dal      Repository
	gateways *dtoInternal.PaymentGateways
}

func New(dal Repository, gateways *dtoInternal.PaymentGateways) *Service {
	return &Service{
		dal:      dal,
		gateways: gateways,
	}
}

func (s *Service) GeneratePaymentURL(ctx context.Context, agg dto.Agg) (orderID string, url string, err error) {
	payment := agg.Payment

	switch payment.Bank {
	case config.Alpha:
		return s.regOrderAlpha(ctx, payment.PaymentID, agg.InternalTransactionUUID, agg.Amount)
	case config.TBank:
		return s.regOrderTbank(ctx, payment.PaymentID, agg.StudentID, agg.InternalTransactionUUID, agg.Amount)
	case config.Tochka:
		return s.regOrderTochka(ctx, payment.PaymentID, agg.InternalTransactionUUID, agg.Amount)
	}

	return "", "", errors.New("Не могу найти ваш банк")
}

func (s *Service) regOrderAlpha(ctx context.Context, paymentID int64, internalTransactionUUID string, amount int) (orderID, url string, _ error) {
	resp, err := s.gateways.Alfa.RegisterOrder(ctx, alfadto.NewOrderRequest(paymentID, internalTransactionUUID, amount))
	if err != nil {
		if resp != nil {
			err = fmt.Errorf("%s: %s", err.Error(), resp.ErrorMessage)
		}
		logger.Error(ctx, "ошибка при создании платежки в альфабанке", err)
		if err = s.dal.DropTransaction(ctx, internalTransactionUUID); err != nil {
			logger.Error(ctx, "ошибка при удалении транзакции при ошибке от альфабанка", err)
			return "", "", err
		}
		return "", "", errors.New("У банка возникли технические неполадки, пожалуйста, попробуйте чуть позже")
	}

	return resp.OrderID, resp.FormURL, nil
}

func (s *Service) regOrderTbank(ctx context.Context, paymentID, studentID int64, internalTransactionUUID string, amount int) (orderID, url string, _ error) {
	phone, err := s.dal.PhoneByStudent(ctx, studentID)
	if err != nil {
		return "", "", err
	}

	orderID, url, err = s.gateways.TBank.InitPayment(ctx, tbankDto.NewInitRequest(paymentID, internalTransactionUUID, int64(amount), phone))
	if err != nil {
		logger.Error(ctx, "ошибка при создании платежки в т банке", err)
		if err = s.dal.DropTransaction(ctx, internalTransactionUUID); err != nil {
			logger.Error(ctx, "ошибка при удалении транзакции при ошибке от т банка", err)
			return "", "", err
		}
		return "", "", errors.New("У банка возникли технические неполадки, пожалуйста, попробуйте чуть позже")
	}
	return orderID, url, nil
}

func (s *Service) regOrderTochka(ctx context.Context, paymentID int64, internalTransactionUUID string, amount int) (orderID, url string, _ error) {
	resp, err := s.gateways.Tochka.InitPayment(ctx, tochkaDto.NewInitRequest(paymentID, int64(amount)))
	if err != nil {
		logger.Error(ctx, "ошибка при создании платежки в точке", err)
		if err = s.dal.DropTransaction(ctx, internalTransactionUUID); err != nil {
			logger.Error(ctx, "ошибка при удалении транзакции при ошибке от точки", err)
			return "", "", err
		}
		return "", "", errors.New("У банка возникли технические неполадки, пожалуйста, попробуйте чуть позже")
	}

	return resp.OperationID, resp.PaymentLink, nil
}
