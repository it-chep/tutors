package dto

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"sort"
	"strconv"

	"github.com/it-chep/tutors.git/internal/config"
	"github.com/samber/lo"
)

const (
	rub = "643"
	//rub = "810"
	ru = "ru"
)

type Credentials struct {
	BaseURL         string
	CredByPaymentID map[int64]config.TBankCred
}

type InitRequest struct {
	// Required
	TerminalKey string `json:"TerminalKey"`
	Amount      int64  `json:"Amount"`
	OrderID     string `json:"OrderId"`
	Token       string `json:"Token"`

	// Optional
	Description     string `json:"Description,omitempty"`
	PayType         string `json:"PayType,omitempty"` // "O" or "T"
	Language        string `json:"Language,omitempty"`
	NotificationURL string `json:"NotificationURL,omitempty"`
	SuccessURL      string `json:"SuccessURL,omitempty"`
	FailURL         string `json:"FailURL,omitempty"`

	paymentID int64
}

func NewInitRequest(paymentID int64, orderID string, amount int64) *InitRequest {
	return &InitRequest{
		Amount:  amount * 100,
		OrderID: orderID,

		PayType:         "O",
		Language:        "ru",
		SuccessURL:      "https://t.me/Payments_A_bot",
		FailURL:         "https://t.me/Payments_A_bot",
		NotificationURL: "https://100rep.ru/callback/tbank",
		Description:     "Оплата консультаций репетитора",
		paymentID:       paymentID,
	}
}

func (r *InitRequest) GenerateToken(password string) {
	data := []struct{ Name, Val string }{
		{Name: "Amount", Val: strconv.FormatInt(r.Amount, 10)},
		{Name: "Description", Val: r.Description},
		{Name: "FailURL", Val: r.FailURL},
		{Name: "Language", Val: r.Language},
		{Name: "NotificationURL", Val: r.NotificationURL},
		{Name: "OrderId", Val: r.OrderID},
		{Name: "PayType", Val: r.PayType},
		{Name: "Password", Val: password},
		{Name: "SuccessURL", Val: r.SuccessURL},
		{Name: "TerminalKey", Val: r.TerminalKey},
	}

	// сортируем по ключу (Name)
	sort.Slice(data, func(i, j int) bool {
		return data[i].Name < data[j].Name
	})

	var str string
	for _, v := range data {
		str += v.Val
	}
	// считаем sha256
	h := sha256.Sum256([]byte(str))
	r.Token = hex.EncodeToString(h[:])
}

func (r *InitRequest) ToHttp(ctx context.Context, cred Credentials) *http.Request {
	r.TerminalKey = cred.CredByPaymentID[r.paymentID].TerminalKey
	r.GenerateToken(cred.CredByPaymentID[r.paymentID].Password)

	body, _ := json.Marshal(r)
	httpReq, _ := http.NewRequestWithContext(ctx, http.MethodPost, cred.BaseURL+"Init", bytes.NewReader(body))
	httpReq.Header.Set("Content-Type", "application/json")
	return httpReq
}

type InitResponse struct {
	Success     bool   `json:"Success"`
	ErrorCode   string `json:"ErrorCode"`
	TerminalKey string `json:"TerminalKey"`
	Status      string `json:"Status"`
	PaymentID   string `json:"PaymentId"`
	OrderID     string `json:"OrderId"`
	Amount      int64  `json:"Amount"`
	PaymentURL  string `json:"PaymentURL"`
	Message     string `json:"Message"`
	Details     string `json:"Details"`
}

type PaymentNotification struct {
	TerminalKey string `json:"TerminalKey"`
	OrderID     string `json:"OrderId"`
	Success     bool   `json:"Success"`
	Status      string `json:"Status"`
	PaymentID   string `json:"PaymentId"`
	ErrorCode   string `json:"ErrorCode"`
	Amount      int64  `json:"Amount"`
	CardID      string `json:"CardId,omitempty"`
	RebillID    string `json:"RebillId,omitempty"`
	Token       string `json:"Token"`
}

type GetOrderRequest struct {
	// Required
	TerminalKey string `json:"TerminalKey"`
	OrderID     string `json:"OrderId"`
	Token       string `json:"Token"`

	paymentID int64
}

func NewGetOrderRequest(paymentID int64, orderID string) *GetOrderRequest {
	return &GetOrderRequest{
		OrderID:   orderID,
		paymentID: paymentID,
	}
}

func (r *GetOrderRequest) GenerateToken(password string) {
	data := []struct{ Name, Val string }{
		{Name: "OrderID", Val: r.OrderID},
		{Name: "Password", Val: password},
		{Name: "TerminalKey", Val: r.TerminalKey},
	}

	// сортируем по ключу (Name)
	sort.Slice(data, func(i, j int) bool {
		return data[i].Name < data[j].Name
	})

	var str string
	for _, v := range data {
		str += v.Val
	}
	// считаем sha256
	h := sha256.Sum256([]byte(str))
	r.Token = hex.EncodeToString(h[:])
}

func (r *GetOrderRequest) ToHttp(ctx context.Context, cred Credentials) *http.Request {
	r.TerminalKey = cred.CredByPaymentID[r.paymentID].TerminalKey
	r.GenerateToken(cred.CredByPaymentID[r.paymentID].Password)

	body, _ := json.Marshal(r)
	httpReq, _ := http.NewRequestWithContext(ctx, http.MethodPost, cred.BaseURL+"CheckOrder", bytes.NewReader(body))
	httpReq.Header.Set("Content-Type", "application/json")
	return httpReq
}

type GetOrderResponse struct {
	TerminalKey string                  `json:"TerminalKey"`
	OrderID     string                  `json:"OrderId"`
	Success     bool                    `json:"Success"`
	ErrorCode   string                  `json:"ErrorCode"`
	Message     string                  `json:"Message"`
	Details     string                  `json:"Details"`
	Payments    []CheckOrderPaymentInfo `json:"Payments"`
}

type CheckOrderPaymentInfo struct {
	Amount  int64  `json:"Amount"`
	Status  string `json:"Status"`
	Success bool   `json:"Success"`
}

func (r *GetOrderResponse) IsPaid() bool {
	if !r.Success {
		return false
	}

	payment := lo.FirstOrEmpty(r.Payments)
	return payment.Success && payment.Status == "CONFIRMED"
}

func (r *GetOrderResponse) Cancelled() bool {
	payment := lo.FirstOrEmpty(r.Payments)
	return !payment.Success && payment.Status == "REJECTED"
}
