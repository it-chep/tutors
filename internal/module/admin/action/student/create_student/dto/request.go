package dto

type CreateRequest struct {
	FirstName  string
	LastName   string
	MiddleName string
	Phone      string
	Tg         string

	CostPerHour string
	SubjectID   int64
	TutorID     int64

	ParentFullName string
	ParentPhone    string
	ParentTg       string

	TgAdminUsername string
}
