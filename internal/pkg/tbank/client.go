package tbank

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/it-chep/tutors.git/internal/pkg/tbank/dto"
)

type Client struct {
	cred       dto.Credentials
	httpClient *http.Client
}

func NewClient(cred dto.Credentials) *Client {
	return &Client{
		cred:       cred,
		httpClient: &http.Client{},
	}
}

func (c *Client) InitPayment(ctx context.Context, req *dto.InitRequest) (orderID, paymentUrl string, _ error) {
	resp, err := c.httpClient.Do(req.ToHttp(ctx, c.cred))
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()

	var result dto.InitResponse
	if err = json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", "", err
	}

	if !result.Success {
		return "", "", errors.New(result.Message)
	}

	qrResp, err := c.httpClient.Do(result.ToQr(ctx, req.PaymentID(), c.cred))
	if err != nil {
		return result.OrderID, result.PaymentURL, nil
	}
	defer qrResp.Body.Close()

	var qrResult dto.GetQrResponse
	if err = json.NewDecoder(qrResp.Body).Decode(&qrResult); err != nil {
		return result.OrderID, result.PaymentURL, nil
	}

	if !qrResult.Success {
		return result.OrderID, result.PaymentURL, nil
	}

	return result.OrderID, qrResult.Data, nil
}

func (c *Client) GetOrderStatus(ctx context.Context, req *dto.GetOrderRequest) (*dto.GetOrderResponse, error) {
	resp, err := c.httpClient.Do(req.ToHttp(ctx, c.cred))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result dto.GetOrderResponse
	if err = json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) KnownTerminal(ctx context.Context, terminal string) bool {
	for _, cred := range c.cred.CredByPaymentID {
		if cred.TerminalKey == terminal {
			return true
		}
	}
	return false
}
