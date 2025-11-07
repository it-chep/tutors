package student_by_id

type Student struct {
	ID         int64  `json:"id"`
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	MiddleName string `json:"middle_name"`
	Phone      string `json:"phone"`
	Tg         string `json:"tg"`

	CostPerHour string `json:"cost_per_hour"`
	SubjectName string `json:"subject_name"`
	TutorID     int64  `json:"tutor_id"`
	TutorName   string `json:"tutor_name"`

	ParentFullName  string `json:"parent_full_name"`
	ParentPhone     string `json:"parent_phone"`
	ParentTg        string `json:"parent_tg"`
	TgAdminUsername string `json:"tg_admin_username"`

	Balance             string `json:"balance"`
	IsOnlyTrialFinished bool   `json:"is_only_trial_finished"`
	IsBalanceNegative   bool   `json:"is_balance_negative"`
	IsNewbie            bool   `json:"is_newbie"`

	ParentTgID int64 `json:"tg_id"`
}

type Response struct {
	Student Student `json:"student"`
}
