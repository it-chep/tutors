package create_tutor

type Request struct {
	FullName    string `json:"full_name"`
	Phone       string `json:"phone"`
	Tg          string `json:"tg"`
	CostPerHour string `json:"cost_per_hour"`
	SubjectID   int64  `json:"subject_id"`
	Email       string `json:"email"`
	AdminID     int64  `json:"admin_id"`
}
