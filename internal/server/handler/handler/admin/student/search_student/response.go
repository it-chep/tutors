package search_student

type Student struct {
	ID         int64  `json:"id"`
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	MiddleName string `json:"middle_name"`

	ParentFullName string `json:"parent_full_name"`
}
type Response struct {
	Students []Student `json:"students"`
}
