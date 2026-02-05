package handler

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/it-chep/tutors.git/internal/module/bot/dto"
	"github.com/it-chep/tutors.git/internal/pkg/logger"
	"net/http"
	"strings"
)

func (h *Handler) bot() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		event, err := h.botParser.HandleUpdate(r)
		if err != nil {
			logger.Error(r.Context(), fmt.Sprintf("[ERROR] Ошибка при хендлинге ивента, TGID: %d", event.SentFrom().ID), err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if event.SentFrom() == nil ||
			event.FromChat() == nil {
			return
		}

		msg := dto.Message{
			User:   event.SentFrom().ID,
			Text:   getTxt(event),
			ChatID: event.FromChat().ID,
		}

		if err = h.botModule.Route(r.Context(), msg); err != nil {
			logger.Error(
				r.Context(),
				fmt.Sprintf(
					"[ERROR] Ошибка при обработке ивента, TGID: %d, username: %s, lastname: %s, firstname: %s",
					event.SentFrom().ID,
					event.SentFrom().UserName,
					event.SentFrom().LastName,
					event.SentFrom().FirstName,
				), err)
			w.WriteHeader(http.StatusBadGateway)
			return
		}
	}
}

func getTxt(update *tgbotapi.Update) string {
	txt := ""

	if update.Message != nil {
		txt = update.Message.Text
		if txt == "/start" {
			return txt
		}
		if strings.Contains(txt, "/start ") {
			return txt[len("/start "):]
		}
	} else if update.CallbackQuery != nil {
		return update.CallbackQuery.Data
	}

	return txt
}
