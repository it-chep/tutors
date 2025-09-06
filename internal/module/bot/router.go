package bot

import (
	"context"

	"github.com/it-chep/tutors.git/internal/module/bot/action/commands/start"
	"github.com/it-chep/tutors.git/internal/module/bot/dto"
)

func (b *Bot) Route(ctx context.Context, msg dto.Message) error {
	switch msg.Text {
	case "/start":
		return b.Actions.Start.Start(msg)
	case start.GetBalance:
		return b.Actions.GetBalance.GetBalance(ctx, msg)
	case start.TopUpBalance:

	}
	return nil
}
