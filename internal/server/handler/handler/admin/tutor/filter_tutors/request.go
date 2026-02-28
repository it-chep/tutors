package filter_tutors

import (
	"github.com/it-chep/tutors.git/internal/module/admin/action/tutor/filter_tutors/dto"
)

type Request struct {
	AdminsUsernamesIDs []int64 `json:"tg_admins_usernames_ids"`
	IsArchived         bool    `json:"is_archive"`
}

func (r Request) ToFilterRequest() dto.FilterRequest {
	return dto.FilterRequest{
		TgUsernameIDs: r.AdminsUsernamesIDs,
		IsArchive:     r.IsArchived,
	}
}
