package create_assistant

type Request struct {
	Email    string `json:"email"`
	FullName string `json:"full_name"`
	Tg       string `json:"tg"`
	Phone    string `json:"phone"`
}
