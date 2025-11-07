package filter_students

type Request struct {
	AdminsUsernames []string `json:"tg_admins_usernames"`
	IsLost          bool     `json:"is_lost"`
}
