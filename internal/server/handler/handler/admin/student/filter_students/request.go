package filter_students

import (
	archive "github.com/it-chep/tutors.git/internal/module/admin/action/student/archive_filter/dto"
	"github.com/it-chep/tutors.git/internal/module/admin/action/student/filter_students/dto"
)

type Request struct {
	AdminsUsernames    []string `json:"tg_admins_usernames"`
	AdminsUsernamesIDs []int64  `json:"tg_admins_usernames_ids"`
	IsLost             bool     `json:"is_lost"`
	IsArchived         bool     `json:"is_archive"`
	PaymentIDs         []int64  `json:"payment_ids"`
}

func (r Request) ToFilterRequest() dto.FilterRequest {
	return dto.FilterRequest{
		IsLost:        r.IsLost,
		TgUsernames:   r.AdminsUsernames,
		TgUsernameIDs: r.AdminsUsernamesIDs,
		PaymentIDs:    r.PaymentIDs,
	}
}

func (r Request) ToArchiveFilterRequest() archive.FilterRequest {
	return archive.FilterRequest{
		IsLost:        r.IsLost,
		TgUsernames:   r.AdminsUsernames,
		TgUsernameIDs: r.AdminsUsernamesIDs,
		PaymentIDs:    r.PaymentIDs,
	}
}
