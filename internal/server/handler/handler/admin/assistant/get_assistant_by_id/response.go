package get_assistant_by_id

type Assistant struct {
	ID           int64    `json:"id"`
	FullName     string   `json:"full_name"`
	Tg           string   `json:"tg"`
	Phone        string   `json:"phone"`
	AvailableTgs []string `json:"tg_admins_usernames"`
}

type Response struct {
	Assistant Assistant `json:"assistant"`
}
