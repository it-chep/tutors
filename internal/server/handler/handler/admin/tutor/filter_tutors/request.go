package filter_tutors

import (
	archiveDto "github.com/it-chep/tutors.git/internal/module/admin/action/tutor/archive_filter/dto"
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

func (r Request) ToArchiveFilterRequest() archiveDto.FilterRequest {
	return archiveDto.FilterRequest{
		TgUsernameIDs: r.AdminsUsernamesIDs,
	}
}
