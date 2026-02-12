package generate_payment_url

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/it-chep/tutors.git/internal/config"
	dtoInternal "github.com/it-chep/tutors.git/internal/dto"
	"github.com/it-chep/tutors.git/internal/module/admin/action/generate_payment_url/dal"
	"github.com/it-chep/tutors.git/internal/module/admin/action/generate_payment_url/dto"
	"github.com/it-chep/tutors.git/internal/module/admin/action/generate_payment_url/service/payment"
	"github.com/it-chep/tutors.git/internal/module/admin/action/generate_payment_url/service/payments_generator"
	"github.com/it-chep/tutors.git/internal/module/admin/action/generate_payment_url/service/validator"
	"github.com/it-chep/tutors.git/internal/pkg/logger"
	"github.com/it-chep/tutors.git/internal/pkg/payment_hash"
	"github.com/jackc/pgx/v5/pgxpool"
	"net/http"
	"strconv"
)

type Action struct {
	dal               *dal.Repository
	paymentService    *payment.Service
	paymentsGenerator *payments_generator.Service
	validator         *validator.Service
}

func NewAction(pool *pgxpool.Pool, gateways *dtoInternal.PaymentGateways, paymentByAdmin config.PaymentsByAdmin) *Action {
	repo := dal.NewRepository(pool)
	return &Action{
		dal:               repo,
		validator:         validator.New(repo),
		paymentService:    payment.NewService(repo, paymentByAdmin),
		paymentsGenerator: payments_generator.New(repo, gateways),
	}
}

func (a *Action) Handle() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		var req dto.Request
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "невалидный запрос", http.StatusBadRequest)
			return
		}

		atoi, err := strconv.Atoi(req.Amount)
		if err != nil {
			logger.Error(ctx, fmt.Sprintf("Невалидное число : %s", req.Amount), err)
			http.Error(w, "невалидный запрос", http.StatusBadRequest)
			return
		}
		amount := int64(atoi)

		studentID, studentUUID, err := payment_hash.DecryptPaymentHash(chi.URLParam(r, "hash"))
		if err != nil {
			logger.Error(ctx, fmt.Sprintf("ошибка дешифрования хэша: %s", chi.URLParam(r, "hash")), err)
			http.Error(w, "невалидный запрос", http.StatusBadRequest)
			return
		}

		err = a.validator.Validate(ctx, studentID, amount, studentUUID)
		if err != nil {
			logger.Error(ctx, fmt.Sprintf("ошибка валидации студента: %d", studentID), err)
			http.Error(w, "невалидный запрос", http.StatusBadRequest)
			return
		}

		paymentData := a.paymentService.GetPayment(ctx, studentID)

		internalTransactionUUID, err := a.dal.InitTransaction(ctx, studentID, amount, paymentData.PaymentID)
		if err != nil {
			logger.Error(ctx, fmt.Sprintf("ошибка создания транзакции studentID: %d", studentID), err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		orderID, url, err := a.paymentsGenerator.GeneratePaymentURL(ctx, dto.Agg{
			InternalTransactionUUID: internalTransactionUUID,
			Payment:                 paymentData,
			Amount:                  int(amount),
		})
		if orderID == "" || err != nil {
			logger.Error(ctx, fmt.Sprintf("ошибка создания заказа studentID: %d, transactionID: %s", studentID, internalTransactionUUID), err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		if err = a.dal.SetTransactionOrder(ctx, internalTransactionUUID, orderID); err != nil {
			logger.Error(ctx, fmt.Sprintf("ошибка обновления транзакции orderID: %s studentID: %d, transactionID: %s", orderID, studentID, internalTransactionUUID), err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err = json.NewEncoder(w).Encode(dto.Response{PaymentURL: url}); err != nil {
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		return
	}
}
