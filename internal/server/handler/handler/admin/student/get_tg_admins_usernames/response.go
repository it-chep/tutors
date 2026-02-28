package get_tg_admins_usernames

type TgAdminUsername struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type Response struct {
	Usernames []TgAdminUsername `json:"tg_admins"`
}
