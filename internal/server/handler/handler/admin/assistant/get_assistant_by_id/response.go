package get_assistant_by_id

type Assistant struct {
	ID       int64  `json:"id"`
	FullName string `json:"full_name"`
	Tg       string `json:"tg"`
	Phone    string `json:"phone"`
}

type Response struct {
	Assistant Assistant `json:"assistant"`
}
