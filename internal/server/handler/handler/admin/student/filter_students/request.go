package filter_students

import (
	archive "github.com/it-chep/tutors.git/internal/module/admin/action/student/archive_filter/dto"
	"github.com/it-chep/tutors.git/internal/module/admin/action/student/filter_students/dto"
)

type Request struct {
	AdminsUsernames []string `json:"tg_admins_usernames"`
	IsLost          bool     `json:"is_lost"`
	IsArchived      bool     `json:"is_archive"`
}

func (r Request) ToFilterRequest() dto.FilterRequest {
	return dto.FilterRequest{
		IsLost:      r.IsLost,
		TgUsernames: r.AdminsUsernames,
	}
}

func (r Request) ToArchiveFilterRequest() archive.FilterRequest {
	return archive.FilterRequest{
		IsLost:      r.IsLost,
		TgUsernames: r.AdminsUsernames,
	}
}
