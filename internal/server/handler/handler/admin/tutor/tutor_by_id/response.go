package tutor_by_id

type Tutor struct {
	ID              int64  `json:"id"`
	FullName        string `json:"full_name"`
	Phone           string `json:"phone"`
	Tg              string `json:"tg"`
	CostPerHour     string `json:"cost_per_hour"`
	SubjectName     string `json:"subject_name"`
	SubjectID       int64  `json:"subject_id"`
	CreatedAt       string `json:"created_at"`
	TgAdminUsername string `json:"tg_admin_username"`
	IsArchive       bool   `json:"is_archive"`
}

type Response struct {
	Tutor Tutor `json:"tutor"`
}
