package create_admin

type Request struct {
	Email    string `json:"email"`
	FullName string `json:"full_name"`
}
