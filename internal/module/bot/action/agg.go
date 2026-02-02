package action

import (
	"github.com/it-chep/tutors.git/internal/config"
	dtoInternal "github.com/it-chep/tutors.git/internal/dto"
	"github.com/it-chep/tutors.git/internal/module/bot/action/commands/acquaintance"
	"github.com/it-chep/tutors.git/internal/module/bot/action/commands/auth_user"
	"github.com/it-chep/tutors.git/internal/module/bot/action/commands/get_balance"
	"github.com/it-chep/tutors.git/internal/module/bot/action/commands/get_lessons"
	"github.com/it-chep/tutors.git/internal/module/bot/action/commands/info"
	"github.com/it-chep/tutors.git/internal/module/bot/action/commands/start"
	"github.com/it-chep/tutors.git/internal/module/bot/action/commands/top_up_balance"
	"github.com/it-chep/tutors.git/internal/pkg/tg_bot"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Agg struct {
	Start        *start.Action
	Acquaintance *acquaintance.Action
	GetBalance   *get_balance.Action
	TopUpBalance *top_up_balance.Action
	GetLessons   *get_lessons.Action
	Info         *info.Action
	AuthUser     *auth_user.Action
}

func NewAgg(pool *pgxpool.Pool, bot *tg_bot.Bot, config *config.Config, gateways *dtoInternal.PaymentGateways) *Agg {
	return &Agg{
		Start:        start.NewAction(pool, bot),
		Acquaintance: acquaintance.NewAction(pool, bot),
		GetBalance:   get_balance.NewAction(pool, bot),
		TopUpBalance: top_up_balance.NewAction(pool, gateways, config.PaymentConfig.PaymentsByAdmin, bot),
		GetLessons:   get_lessons.NewAction(pool, bot),
		Info:         info.New(bot),
		AuthUser:     auth_user.NewAction(pool, bot),
	}
}
