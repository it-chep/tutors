package get_user_info

import (
	"encoding/json"
	"net/http"

	get_user_info_dal "github.com/it-chep/tutors.git/internal/module/admin/action/auth/get_user_info/dal"
	userCtx "github.com/it-chep/tutors.git/pkg/context"
	"github.com/jackc/pgx/v5/pgxpool"
)

type User struct {
	ID   int64  `json:"id"`
	Role string `json:"role"`
}

type Response struct {
	User User `json:"user"`
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

		user, err := a.dal.GetUser(ctx, userCtx.UserFromContext(ctx))
		if err != nil {
			http.Error(w, "Что-то пошло не так, попробуйте позже", http.StatusInternalServerError)
			return
		}

		resp := Response{
			User: User{
				ID:   user.ID,
				Role: user.Role.FrontString(),
			},
		}
		w.Header().Set("Content-Type", "application/json")
		if err = json.NewEncoder(w).Encode(resp); err != nil {
			http.Error(w, "Что-то пошло не так, попробуйте позже", http.StatusInternalServerError)
			return
		}
	}
}
