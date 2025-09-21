package get_admin_by_id

type Admin struct {
	ID       int64  `json:"id"`
	FullName string `json:"full_name"`
	Tg       string `json:"tg"`
	Phone    string `json:"phone"`
}

type Response struct {
	Admin Admin `json:"admin"`
}
