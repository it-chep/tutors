package alpha

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"net/http"

	alpha_dal "github.com/it-chep/tutors.git/internal/module/admin/alpha/dal"
	"github.com/samber/lo"
)

type WebHookAlpha struct {
	dal *alpha_dal.Repository
}

func NewWebHookAlpha(dal *alpha_dal.Repository) *WebHookAlpha {
	return &WebHookAlpha{dal: dal}
}

func (hook *WebHookAlpha) Handle() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}
		defer func() {
			_ = r.Body.Close()
		}()

		secretKey := "your-webhook-secret-key"
		signature := r.Header.Get("X-Signature-SHA256")

		if !verifyWebhookSignature(body, signature, secretKey) {
			http.Error(w, "Invalid signature", http.StatusUnauthorized)
			return
		}

		var webhook WebhookRequest
		if err = json.Unmarshal(body, &webhook); err != nil {
			http.Error(w, "Invalid JSON format", http.StatusBadRequest)
			return
		}

		if err = webhook.Validate(); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err = hook.processWebhook(r.Context(), webhook); err != nil {
			http.Error(w, err.Error(), http.StatusBadGateway)
		}

		w.WriteHeader(http.StatusOK)
	}
}

func verifyWebhookSignature(body []byte, signature, secret string) bool {
	if signature == "" {
		return false
	}

	h := hmac.New(sha256.New, []byte(secret))
	h.Write(body)
	expectedSign := hex.EncodeToString(h.Sum(nil))

	return hmac.Equal([]byte(expectedSign), []byte(signature))
}

func (hook *WebHookAlpha) processWebhook(ctx context.Context, webhook WebhookRequest) error {
	if webhook.EventType != "ORDER_STATUS_UPDATED" {
		return nil
	}
	if !lo.Contains([]string{"APPROVED", "CONFIRMED"}, webhook.Data.Order.Status) {
		return nil
	}

	return hook.dal.UpdateBalance(ctx, webhook.Data.Order.OrderNumber)
}
