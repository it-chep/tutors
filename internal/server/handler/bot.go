package handler

import (
	"github.com/it-chep/tutors.git/internal/module/bot/dto"
	"github.com/it-chep/tutors.git/internal/pkg/logger"
	"net/http"
)

func (h *Handler) bot() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		event, err := h.botParser.HandleUpdate(r)
		if err != nil {
			logger.Error(r.Context(), "[ERROR] Ошибка при хендлинге ивента", err)
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
			logger.Error(r.Context(), "[ERROR] Ошибка при обработке ивента", err)
			w.WriteHeader(http.StatusBadGateway)
			return
		}
	}
}
