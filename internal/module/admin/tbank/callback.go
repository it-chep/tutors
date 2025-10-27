package tbank

import (
	"io"
	"net/http"

	"github.com/it-chep/tutors.git/internal/module/admin/job/order_checker"
	"github.com/it-chep/tutors.git/internal/pkg/logger"
)

type CallbackTbank struct {
	checker *order_checker.TransactionChecker
}

func NewCallbackTbank(checker *order_checker.TransactionChecker) *CallbackTbank {
	return &CallbackTbank{checker: checker}
}

func (hook *CallbackTbank) Handle() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Message(r.Context(), "получен колбек от тбанка")

		callBack, err := hook.webHook(r)
		if err != nil {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}

		if callBack.Status != "CONFIRMED" {
			return
		}
		if err = hook.checker.ConfirmOrder(r.Context(), callBack.OrderID, callBack.TerminalKey); err != nil {
			http.Error(w, err.Error(), http.StatusBadGateway)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

func (hook *CallbackTbank) webHook(r *http.Request) (*Callback, error) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = r.Body.Close()
	}()

	webhook := &Callback{}
	return webhook, webhook.UnmarshalJSON(body)
}
