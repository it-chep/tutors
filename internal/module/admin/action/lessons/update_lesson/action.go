package update_lesson

import (
	"context"

	"github.com/it-chep/tutors.git/internal/module/admin/action/lessons/update_lesson/dal"
	"github.com/it-chep/tutors.git/internal/module/admin/action/lessons/update_lesson/dto"
	"github.com/it-chep/tutors.git/internal/pkg/tg_bot"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Action провести обычное занятие
type Action struct {
	dal *dal.Repository
	bot *tg_bot.Bot
}

func New(pool *pgxpool.Pool, bot *tg_bot.Bot) *Action {
	return &Action{
		dal: dal.NewRepository(pool),
		bot: bot,
	}
}

func (a *Action) Do(ctx context.Context, lessonID int64, upd dto.UpdateLesson) error {
	return nil
}
