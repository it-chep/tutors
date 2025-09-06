package student_by_id

type Student struct {
	ID         int64  `json:"id"`
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

	Balance             string `json:"balance"`
	HasButtons          bool   `json:"has_buttons"`
	IsOnlyTrialFinished bool   `json:"is_only_trial_finished"`
	IsBalanceNegative   bool   `json:"is_balance_negative"`
	IsNewbie            bool   `json:"is_newbie"`
}

type Response struct {
	Student Student `json:"student"`
}
