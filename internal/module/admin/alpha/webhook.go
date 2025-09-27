package alpha

import (
	"context"
	"io"
	"net/http"
	"strings"

	alpha_dal "github.com/it-chep/tutors.git/internal/module/admin/alpha/dal"
)

type WebHookAlpha struct {
	dal    *alpha_dal.Repository
	secret string
}

func NewWebHookAlpha(dal *alpha_dal.Repository, secret string) *WebHookAlpha {
	return &WebHookAlpha{dal: dal, secret: secret}
}

func (hook *WebHookAlpha) Handle() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if !hook.checkBearer(r) {
			http.Error(w, "Invalid auth", http.StatusUnauthorized)
			return
		}

		webhook, err := hook.webHook(r)
		if err != nil {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}

		if err = hook.processWebhook(r.Context(), webhook); err != nil {
			http.Error(w, err.Error(), http.StatusBadGateway)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

func (hook *WebHookAlpha) checkBearer(r *http.Request) bool {
	auth := r.Header.Get("Authorization")
	if !strings.HasPrefix(auth, "Bearer ") {
		return false
	}
	return strings.TrimPrefix(auth, "Bearer ") == hook.secret
}

func (hook *WebHookAlpha) webHook(r *http.Request) (*WebhookEnvelope, error) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = r.Body.Close()
	}()

	webhook := &WebhookEnvelope{}
	return webhook, webhook.UnmarshalJSON(body)
}

func (hook *WebHookAlpha) processWebhook(ctx context.Context, webhook *WebhookEnvelope) error {
	return hook.dal.UpdateBalance(ctx, webhook.OrderNumber())
}
