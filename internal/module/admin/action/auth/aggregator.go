package auth

import (
	"github.com/it-chep/tutors.git/internal/config"
	"github.com/it-chep/tutors.git/internal/module/admin/action/auth/check_path_permission"
	"github.com/it-chep/tutors.git/internal/module/admin/action/auth/get_user_info"
	"github.com/it-chep/tutors.git/internal/module/admin/action/auth/login"
	"github.com/it-chep/tutors.git/internal/module/admin/action/auth/logout"
	"github.com/it-chep/tutors.git/internal/module/admin/action/auth/refresh"
	"github.com/it-chep/tutors.git/internal/module/admin/action/auth/register"
	"github.com/it-chep/tutors.git/pkg/smtp"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Aggregator struct {
	// AUTH
	CheckPathPermission *check_path_permission.Action
	Login               *login.Action
	Refresh             *refresh.Action
	Register            *register.Action
	GetUserInfo         *get_user_info.Action
	Logout              *logout.Action
}

func NewAggregator(pool *pgxpool.Pool, smtp *smtp.ClientSmtp, config config.JwtConfig) *Aggregator {
	return &Aggregator{
		CheckPathPermission: check_path_permission.New(pool, config),
		Login:               login.New(pool, smtp, config),
		Refresh:             refresh.New(config),
		Register:            register.New(pool, smtp, config),
		GetUserInfo:         get_user_info.New(pool),
		Logout:              logout.New(),
	}
}
