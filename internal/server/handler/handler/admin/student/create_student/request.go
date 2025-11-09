package create_student

type Request struct {
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	MiddleName string `json:"middle_name"`
	Phone      string `json:"phone"`
	Tg         string `json:"tg"`

	CostPerHour string `json:"cost_per_hour"`
	SubjectID   int64  `json:"subject_id"`
	TutorID     int64  `json:"tutor_id"`

	ParentFullName string `json:"parent_full_name"`
	ParentPhone    string `json:"parent_phone"`
	ParentTg       string `json:"parent_tg"`

	TgAdminUsername *string `json:"tg_admin_username,omitempty"`
}
