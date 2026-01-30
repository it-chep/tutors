package alpha

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/it-chep/tutors.git/internal/module/admin/job/order_checker"
	"github.com/it-chep/tutors.git/internal/pkg/logger"
)

type WebHookAlpha struct {
	checker *order_checker.TransactionChecker
	secret  string
}

func NewWebHookAlpha(checker *order_checker.TransactionChecker, secret string) *WebHookAlpha {
	return &WebHookAlpha{checker: checker, secret: secret}
}

func (hook *WebHookAlpha) Handle() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Message(r.Context(), "получен вебхук от альфы")
		if !hook.checkBearer(r) {
			http.Error(w, "Invalid auth", http.StatusUnauthorized)
			return
		}

		webhook, err := hook.webHook(r)
		logger.Message(r.Context(), fmt.Sprintf("вебхук: %s", string(webhook.Data)))
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
	//auth := r.Header.Get("Authorization")
	//if !strings.HasPrefix(auth, "Bearer ") {
	//	return false
	//}
	//return strings.TrimPrefix(auth, "Bearer ") == hook.secret
	return true
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
	amount := webhook.Amount()
	logger.Message(ctx, fmt.Sprintf("величина поступления на счет: %s", amount.String()))
	//if amount.IsZero() {
	return nil
	//}
	//return hook.checker.UpdateTransactionsByAmount(ctx, amount)
}
