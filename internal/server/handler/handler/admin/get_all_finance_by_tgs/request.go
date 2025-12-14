package get_all_finance_by_tgs

type Request struct {
	From    string `json:"from"`
	To      string `json:"to"`
	AdminID int64  `json:"admin_id"`

	TgUsernames []string `json:"tg_admin_usernames"`
}
