package dto

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

const (
	rub = "810"
	ru  = "ru"
)

type Credentials struct {
	BaseURL  string
	UserName string
	Password string
}

type OrderRequest struct {
	orderNumber string
	amount      int
	currency    string
	returnURL   string
	description string
	language    string
}

func NewOrderRequest(orderNumber string, amount int) *OrderRequest {
	return &OrderRequest{
		orderNumber: orderNumber,
		amount:      amount,
		currency:    rub,
		returnURL:   fmt.Sprintf("https://t.me/Payments_A_bot"),
		description: "Оплата консультаций репетитора",
		language:    ru,
	}
}

func (r OrderRequest) FormData(ctx context.Context, cred Credentials) (*http.Request, error) {
	formData := url.Values{}
	formData.Set("userName", cred.UserName)
	formData.Set("password", cred.Password)
	formData.Set("orderNumber", r.orderNumber)
	formData.Set("amount", fmt.Sprintf("%d", r.amount))
	formData.Set("returnUrl", r.returnURL)
	formData.Set("description", r.description)
	formData.Set("currency", r.currency) // RUB
	formData.Set("language", r.language)

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
	orderID     string
	orderNumber string
}

func NewStatusRequest(orderID, orderNumber string) *StatusRequest {
	return &StatusRequest{
		orderID:     orderID,
		orderNumber: orderNumber,
	}
}

func (r StatusRequest) FormData(ctx context.Context, cred Credentials) (*http.Request, error) {
	formData := url.Values{}
	formData.Set("userName", cred.UserName)
	formData.Set("password", cred.Password)
	formData.Set("orderId", r.orderID)
	formData.Set("orderNumber", r.orderNumber)

	req, err := http.NewRequestWithContext(ctx, "POST", cred.BaseURL+"/getOrderStatus.do", strings.NewReader(formData.Encode()))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	return req, nil
}

type StatusResponse struct {
	OrderNumber  string `json:"orderNumber"`
	OrderStatus  int    `json:"orderStatus"`
	ErrorCode    string `json:"errorCode"`
	ErrorMessage string `json:"errorMessage"`
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
