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
	//path := "/Users/duzyanov/govno/tutors/internal/pkg/alpha/cert/"
	//cert, err := tls.LoadX509KeyPair(path+"client_fullchain.pem", path+"sandbox_key_2026_nopass.key")
	//if err != nil {
	//	log.Fatal("client cert load error:", err)
	//}
	//
	//caCert, err := ioutil.ReadFile(path + "alfa_root_bundle.pem")
	//if err != nil {
	//	log.Fatalf("failed to read CA bundle: %v", err)
	//}
	//caPool := x509.NewCertPool()
	//if !caPool.AppendCertsFromPEM(caCert) {
	//	log.Fatal("failed to append CA bundle")
	//}
	//
	//tr := &http.Transport{
	//	TLSClientConfig: &tls.Config{
	//		Certificates:       []tls.Certificate{cert},
	//		RootCAs:            caPool,
	//		MinVersion:         tls.VersionTLS12,
	//		InsecureSkipVerify: true,
	//	},
	//}

	return &Client{
		cred: cred,
		httpClient: &http.Client{
			//Transport: tr,
			Timeout: 30 * time.Second,
		},
	}
}

func (c *Client) RegisterOrder(ctx context.Context, request *dto.OrderRequest) (*dto.OrderResponse, error) {
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
