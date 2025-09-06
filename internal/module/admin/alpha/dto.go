package alpha

import (
	"fmt"
	"time"
)

type WebhookRequest struct {
	EventID      string    `json:"eventId"`
	EventType    string    `json:"eventType"`
	EventVersion string    `json:"eventVersion"`
	CreatedAt    time.Time `json:"createdAt"`
	Data         struct {
		Order struct {
			OrderID     string `json:"orderId"`
			OrderNumber string `json:"orderNumber"`
			Status      string `json:"status"`
		} `json:"order"`
		Payment struct {
			PaymentID string `json:"paymentId"`
			Amount    struct {
				Value    string `json:"value"`
				Currency string `json:"currency"`
			} `json:"amount"`
			PaymentMethod struct {
				Type string `json:"type"`
				Card struct {
					Last4 string `json:"last4"`
					Brand string `json:"brand"`
				} `json:"card,omitempty"`
			} `json:"paymentMethod"`
		} `json:"payment"`
	} `json:"data"`
}

func (r *WebhookRequest) Validate() error {
	if r.EventID == "" {
		return fmt.Errorf("eventId is required")
	}
	if r.EventType == "" {
		return fmt.Errorf("eventType is required")
	}
	if r.Data.Order.OrderNumber == "" {
		return fmt.Errorf("orderNumber is required")
	}
	if r.Data.Order.Status == "" {
		return fmt.Errorf("order status is required")
	}

	return nil
}

type WebhookResponse struct {
	EventID   string `json:"eventId"`
	Status    string `json:"status"`
	Processed bool   `json:"processed"`
}
