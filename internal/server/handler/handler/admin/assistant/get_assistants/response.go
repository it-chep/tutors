package get_assistants

type Assistant struct {
	ID       int64  `json:"id"`
	FullName string `json:"full_name"`
	Tg       string `json:"tg"`
}

type Response struct {
	Assistants []Assistant `json:"assistants"`
}
