package update_student

import "github.com/it-chep/tutors.git/internal/module/admin/action/student/update_student/dto"

type Request struct {
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	MiddleName string `json:"middle_name"`
	Phone      string `json:"phone"`
	Tg         string `json:"tg"`

	CostPerHour string `json:"cost_per_hour"`

	ParentFullName  string `json:"parent_full_name"`
	ParentPhone     string `json:"parent_phone"`
	ParentTg        string `json:"parent_tg"`
	TgAdminUsername string `json:"tg_admin_username"`
}

func (r Request) ToDto() dto.UpdateRequest {
	return dto.UpdateRequest{
		FirstName:       r.FirstName,
		LastName:        r.LastName,
		MiddleName:      r.MiddleName,
		Phone:           r.Phone,
		Tg:              r.Tg,
		CostPerHour:     r.CostPerHour,
		ParentFullName:  r.ParentFullName,
		ParentPhone:     r.ParentPhone,
		ParentTg:        r.ParentTg,
		TgAdminUsername: r.TgAdminUsername,
	}
}
