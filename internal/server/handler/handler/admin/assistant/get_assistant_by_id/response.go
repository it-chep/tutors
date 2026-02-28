package get_assistant_by_id

type Assistant struct {
	ID           int64             `json:"id"`
	FullName     string            `json:"full_name"`
	Tg           string            `json:"tg"`
	Phone        string            `json:"phone"`
	AvailableTgs []TgAdminUsername `json:"tg_admins_usernames"`
	CreatedAt    string            `json:"created_at"`
}

type TgAdminUsername struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type Response struct {
	Assistant Assistant `json:"assistant"`
}
