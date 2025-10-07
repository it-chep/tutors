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

		if event.SentFrom() == nil ||
			event.FromChat() == nil {
			return
		}

		txt := ""
		if event.Message != nil {
			txt = event.Message.Text
		} else if event.CallbackQuery != nil {
			txt = event.CallbackQuery.Data
		}

		msg := dto.Message{
			User:   event.SentFrom().ID,
			Text:   txt,
			ChatID: event.FromChat().ID,
		}

		if err = h.botModule.Route(r.Context(), msg); err != nil {
			w.WriteHeader(http.StatusBadGateway)
			return
		}
	}
}
