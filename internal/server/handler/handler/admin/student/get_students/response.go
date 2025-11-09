package get_students

type Student struct {
	ID                  int64  `json:"id"`
	FirstName           string `json:"first_name"`
	LastName            string `json:"last_name"`
	MiddleName          string `json:"middle_name"`
	ParentFullName      string `json:"parent_full_name"`
	Tg                  string `json:"tg"`
	IsOnlyTrialFinished bool   `json:"is_only_trial_finished"`
	IsBalanceNegative   bool   `json:"is_balance_negative"`
	IsNewbie            bool   `json:"is_newbie"`
	Balance             string `json:"balance"`
}

type Response struct {
	Students      []Student `json:"students"`
	StudentsCount int64     `json:"students_count"`
}
