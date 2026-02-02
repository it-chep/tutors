package dto

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"

	"github.com/it-chep/tutors.git/internal/config"
)

type Credentials struct {
	BaseURL         string
	CredByPaymentID map[int64]config.TochkaCred
}

type InitReqData struct {
	CustomerCode string   `json:"customerCode"`
	Amount       int64    `json:"amount"`
	Purpose      string   `json:"purpose"`
	PaymentMode  []string `json:"paymentMode"` // "card" | "sbp" | "tinkoff"
	RedirectURL  string   `json:"redirectUrl"`
	SaveCard     bool     `json:"saveCard"`
	paymentID    int64
}

type InitRequest struct {
	Data InitReqData `json:"Data"`
}

func NewInitRequest(paymentID int64, amount int64) *InitRequest {
	return &InitRequest{
		Data: InitReqData{
			Amount:      amount,
			Purpose:     "Оплата консультаций репетитора",
			PaymentMode: []string{"card", "sbp", "tinkoff"},
			RedirectURL: "https://t.me/Payments_A_bot",
			SaveCard:    true,
			paymentID:   paymentID,
		},
	}
}

func (r *InitRequest) ToHttp(ctx context.Context, cred Credentials) (*http.Request, error) {
	r.Data.CustomerCode = cred.CredByPaymentID[r.Data.paymentID].CustomerCode

	body, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		cred.BaseURL+"payments",
		bytes.NewReader(body),
	)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+cred.CredByPaymentID[r.Data.paymentID].JWT)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	return req, nil
}

type InitResponse struct {
	PaymentLink string `json:"paymentLink"`
	OperationID string `json:"operationId"`
}

type InitResponseHTTP struct {
	Data InitResponse `json:"Data"`
}

type GetOrderRequest struct {
	OperationID string
	paymentID   int64
}

func NewGetOrderRequest(paymentID int64, operationID string) *GetOrderRequest {
	return &GetOrderRequest{
		OperationID: operationID,
		paymentID:   paymentID,
	}
}

func (r *GetOrderRequest) ToHttp(ctx context.Context, cred Credentials) (*http.Request, error) {
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		cred.BaseURL+"payments/"+r.OperationID,
		nil,
	)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+cred.CredByPaymentID[r.paymentID].JWT)
	return req, nil
}

type GetOrderResponse struct {
	Status string `json:"status"`
}
type GetOrderResponseOperation struct {
	GetOrderResponse []GetOrderResponse `json:"Operation"`
}

type GetOrderResponseHTTP struct {
	Data GetOrderResponseOperation `json:"Data"`
}

func (r *GetOrderResponse) IsPaid() bool {
	return r.Status == "APPROVED"
}

func (r *GetOrderResponse) Expired() bool {
	return r.Status == "EXPIRED"
}
