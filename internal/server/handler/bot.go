package handler

import (
	"github.com/it-chep/tutors.git/internal/module/bot/dto"
	"net/http"
)

func (h *Handler) bot() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		event, err := h.botParser.HandleUpdate(r)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		msg := dto.Message{
			User:   event.SentFrom().ID,
			ChatID: event.FromChat().ID,
		}
		if event.Message != nil {
			msg.Text = event.Message.Text
		} else if event.CallbackQuery != nil {
			msg.Text = event.CallbackQuery.Data
		} else {
			return
		}
		if err = h.botModule.Route(r.Context(), msg); err != nil {
			w.WriteHeader(http.StatusBadGateway)
			return
		}
	}
}
