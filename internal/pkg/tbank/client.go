package tbank

import (
	"context"
	"encoding/json"
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

func (c *Client) InitPayment(ctx context.Context, req *dto.InitRequest) (*dto.InitResponse, error) {
	resp, err := c.httpClient.Do(req.ToHttp(ctx, c.cred))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result dto.InitResponse
	if err = json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return &result, nil
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
