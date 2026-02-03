package dto

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/it-chep/tutors.git/internal/config"
)

const (
	rub = "643"
	//rub = "810"
	ru = "ru"
)

type Credentials struct {
	BaseURL         string
	CredByPaymentID map[int64]config.AlphaCred
}

type OrderRequest struct {
	orderNumber string
	amount      int
	currency    string
	returnURL   string
	description string
	language    string

	paymentID int64
}

func NewOrderRequest(paymentID int64, orderNumber string, amount int) *OrderRequest {
	return &OrderRequest{
		orderNumber: orderNumber,
		amount:      amount,
		currency:    rub,
		returnURL:   fmt.Sprintf("https://t.me/Payments_A_bot"),
		description: "Оплата консультаций репетитора",
		language:    ru,
		paymentID:   paymentID,
	}
}

func (r OrderRequest) FormData(ctx context.Context, cred Credentials) (*http.Request, error) {
	formData := url.Values{}
	formData.Set("userName", cred.CredByPaymentID[r.paymentID].User)
	formData.Set("password", cred.CredByPaymentID[r.paymentID].Password)
	formData.Set("orderNumber", r.orderNumber)
	formData.Set("amount", fmt.Sprintf("%d", r.amount*100))
	formData.Set("returnUrl", r.returnURL)
	formData.Set("description", r.description)
	formData.Set("currency", r.currency) // RUB
	formData.Set("language", r.language)
	formData.Set("sessionTimeoutSecs", "1200")

	req, err := http.NewRequestWithContext(ctx, "POST", cred.BaseURL+"/register.do", strings.NewReader(formData.Encode()))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	return req, nil
}

type OrderResponse struct {
	OrderID      string `json:"orderId"`
	FormURL      string `json:"formUrl"`
	ErrorCode    string `json:"errorCode"`
	ErrorMessage string `json:"errorMessage"`
}

func (r *OrderResponse) FromHttp(reader io.Reader) error {
	body, err := io.ReadAll(reader)
	if err != nil {
		return err
	}

	if err = json.Unmarshal(body, r); err != nil {
		return err
	}

	if r.ErrorCode != "" {
		return fmt.Errorf("API error: %s (code: %s)", r.ErrorMessage, r.ErrorCode)
	}

	if r.FormURL == "" {
		return fmt.Errorf("касса не сформировала ссылку на оплату")
	}

	return nil
}

type StatusRequest struct {
	orderID   string
	paymentID int64
}

func NewStatusRequest(paymentID int64, orderID string) *StatusRequest {
	return &StatusRequest{
		orderID:   orderID,
		paymentID: paymentID,
	}
}

func (r StatusRequest) FormData(ctx context.Context, cred Credentials) (*http.Request, error) {
	formData := url.Values{}
	formData.Set("userName", cred.CredByPaymentID[r.paymentID].User)
	formData.Set("password", cred.CredByPaymentID[r.paymentID].Password)
	formData.Set("orderId", r.orderID)

	req, err := http.NewRequestWithContext(ctx, "POST", cred.BaseURL+"/getOrderStatus.do", strings.NewReader(formData.Encode()))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	return req, nil
}

type OrderStatus int

func (r OrderStatus) Confirmed() bool {
	return r == 2
}

func (r OrderStatus) Cancelled() bool {
	return r == 3 || r == 6
}

type StatusResponse struct {
	OrderNumber  string      `json:"OrderNumber"`
	OrderStatus  OrderStatus `json:"OrderStatus"`
	ErrorCode    string      `json:"ErrorCode"`
	ErrorMessage string      `json:"ErrorMessage"`
}

func (r *StatusResponse) FromHttp(reader io.Reader) error {
	body, err := io.ReadAll(reader)
	if err != nil {
		return err
	}

	if err = json.Unmarshal(body, r); err != nil {
		return err
	}

	if r.ErrorCode != "" {
		return fmt.Errorf("API error: %s (code: %s)", r.ErrorMessage, r.ErrorCode)
	}

	return nil
}
