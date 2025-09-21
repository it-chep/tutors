package get_admins

type Admin struct {
	ID       int64  `json:"id"`
	FullName string `json:"full_name"`
	Tg       string `json:"tg"`
}

type Response struct {
	Admins []Admin `json:"admins"`
}
