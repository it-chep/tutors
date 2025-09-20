package bot

import (
	"context"

	"github.com/it-chep/tutors.git/internal/module/bot/action/commands/start"
	"github.com/it-chep/tutors.git/internal/module/bot/dto"
)

func (b *Bot) Route(ctx context.Context, msg dto.Message) error {
	switch msg.Text {
	case "/start":
		return b.Actions.Start.Start(ctx, msg)
	case start.GetBalance:
		return b.Actions.GetBalance.GetBalance(ctx, msg)
	case start.TopUpBalance:
		return b.Actions.TopUpBalance.InitTransaction(ctx, msg)
	default:
		if ok, err := b.Actions.Acquaintance.OnRegistration(ctx, msg); err == nil && ok {
			return b.Actions.Acquaintance.MakeKnown(ctx, msg)
		}

		if b.Actions.TopUpBalance.TransactionExists(ctx, msg) {
			return b.Actions.TopUpBalance.SetAmount(ctx, msg)
		}
	}
	return nil
}
