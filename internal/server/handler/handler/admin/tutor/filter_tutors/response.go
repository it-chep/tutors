package filter_tutors

type Tutor struct {
	ID                 int64  `json:"id"`
	FullName           string `json:"full_name"`
	Tg                 string `json:"tg"`
	HasBalanceNegative bool   `json:"has_balance_negative"`
	HasOnlyTrial       bool   `json:"has_only_trial"`
	HasNewBie          bool   `json:"has_newbie"`
	IsFailer           bool   `json:"is_failer"`
	TgAdminUsername    string `json:"tg_admin_username"`
}

type Response struct {
	Tutors      []Tutor `json:"tutors"`
	TutorsCount int64   `json:"tutors_count"`
}
