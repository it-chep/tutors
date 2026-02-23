package filter_tutors

import (
	"github.com/it-chep/tutors.git/internal/module/admin/action/tutor/filter_tutors/dto"
)

type Request struct {
	AdminsUsernames []string `json:"tg_admins_usernames"`
}

func (r Request) ToFilterRequest() dto.FilterRequest {
	return dto.FilterRequest{
		TgUsernames: r.AdminsUsernames,
	}
}
