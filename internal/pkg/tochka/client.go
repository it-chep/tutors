package tochka

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/it-chep/tutors.git/internal/pkg/tochka/dto"
	"github.com/samber/lo"
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
	httpReq, err := req.ToHttp(ctx, c.cred)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("init payment failed: %s", resp.Status)
	}

	defer resp.Body.Close()

	var result dto.InitResponseHTTP
	if err = json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result.Data, nil
}

func (c *Client) GetOrderStatus(ctx context.Context, req *dto.GetOrderRequest) (*dto.GetOrderResponse, error) {
	httpReq, err := req.ToHttp(ctx, c.cred)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("tochka GetOrderStatus failed: %s", resp.Status)
	}
	defer resp.Body.Close()

	var result dto.GetOrderResponseHTTP
	if err = json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	if len(result.Data.GetOrderResponse) != 1 {
		return nil, fmt.Errorf("получение заказа с точки, ожидали 1, получили %d", len(result.Data.GetOrderResponse))
	}

	response := lo.FirstOrEmpty(result.Data.GetOrderResponse)
	return &response, nil
}
