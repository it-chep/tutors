package get_user_info

import (
	"encoding/json"
	"github.com/it-chep/tutors.git/internal/pkg/logger"
	"net/http"

	get_user_info_dal "github.com/it-chep/tutors.git/internal/module/admin/action/auth/get_user_info/dal"
	userCtx "github.com/it-chep/tutors.git/pkg/context"
	"github.com/jackc/pgx/v5/pgxpool"
)

type User struct {
	ID      int64  `json:"id"`
	Role    string `json:"role"`
	TutorID int64  `json:"tutor_id"`
}

type Response struct {
	User          User            `json:"user"`
	PaidFunctions map[string]bool `json:"paid_functions"`
}

type Action struct {
	dal *get_user_info_dal.Repository
}

func New(pool *pgxpool.Pool) *Action {
	return &Action{
		dal: get_user_info_dal.NewRepository(pool),
	}
}

func (a *Action) Handle() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		user, err := a.dal.GetUser(ctx, userCtx.UserIDFromContext(ctx))
		if err != nil {
			http.Error(w, "Что-то пошло не так, попробуйте позже", http.StatusInternalServerError)
			logger.Error(r.Context(), "get userinfo Ошибка при юзере", err)
			return
		}

		paid, err := a.dal.GetPaidFunctions(ctx, user.AdminID)
		if err != nil {
			http.Error(w, "Ошибка из GetPaidFunctions", http.StatusInternalServerError)
			logger.Error(r.Context(), "Ошибка при GetPaidFunctions", err)
			return
		}

		functions := make(map[string]bool)
		if paid != nil {
			functions = paid.PaidFunctions
		}

		resp := Response{
			User: User{
				ID:      user.ID,
				Role:    user.Role.FrontString(),
				TutorID: user.TutorID,
			},
			PaidFunctions: functions,
		}
		w.Header().Set("Content-Type", "application/json")
		if err = json.NewEncoder(w).Encode(resp); err != nil {
			http.Error(w, "Что-то пошло не так, попробуйте позже", http.StatusInternalServerError)
			return
		}
	}
}
