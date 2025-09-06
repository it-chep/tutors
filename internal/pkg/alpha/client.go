package alfa

import (
	"context"
	"net/http"
	"time"

	"github.com/it-chep/tutors.git/internal/pkg/alpha/dto"
)

type Client struct {
	cred       dto.Credentials
	httpClient *http.Client
}

func NewClient(cred dto.Credentials) *Client {
	return &Client{
		cred: cred,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (c *Client) RegisterOrder(ctx context.Context, request dto.OrderRequest) (*dto.OrderResponse, error) {
	req, err := request.FormData(ctx, c.cred)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	response := &dto.OrderResponse{}
	return response, response.FromHttp(resp.Body)
}

func (c *Client) GetOrderStatus(ctx context.Context, request dto.StatusRequest) (*dto.StatusResponse, error) {
	req, err := request.FormData(ctx, c.cred)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	response := &dto.StatusResponse{}
	return response, response.FromHttp(resp.Body)
}
