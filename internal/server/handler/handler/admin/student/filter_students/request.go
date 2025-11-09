package filter_students

import "github.com/it-chep/tutors.git/internal/module/admin/action/student/filter_students/dto"

type Request struct {
	AdminsUsernames []string `json:"tg_admins_usernames"`
	IsLost          bool     `json:"is_lost"`
}

func (r Request) ToFilterRequest() dto.FilterRequest {
	return dto.FilterRequest{
		IsLost:      r.IsLost,
		TgUsernames: r.AdminsUsernames,
	}
}
