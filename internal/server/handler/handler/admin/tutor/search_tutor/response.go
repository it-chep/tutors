package search_tutor

type Tutor struct {
	ID       int64  `json:"id"`
	FullName string `json:"full_name"`
}
type Response struct {
	Tutors []Tutor `json:"tutors"`
}
