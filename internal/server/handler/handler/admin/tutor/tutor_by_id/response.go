package tutor_by_id

type Tutor struct {
	ID          int64  `json:"id"`
	FullName    string `json:"full_name"`
	Phone       string `json:"phone"`
	Tg          string `json:"tg"`
	CostPerHour string `json:"cost_per_hour"`
	SubjectID   int64  `json:"subject_id"`
}

type Response struct {
	Tutor Tutor `json:"tutor"`
}
