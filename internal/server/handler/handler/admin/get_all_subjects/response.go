package get_all_subjects

type Subject struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type Response struct {
	Subjects []Subject `json:"subjects"`
}
