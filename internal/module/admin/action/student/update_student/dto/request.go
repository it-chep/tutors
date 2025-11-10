package dto

type UpdateRequest struct {
	FirstName  string
	LastName   string
	MiddleName string
	Phone      string
	Tg         string

	CostPerHour string

	ParentFullName string
	ParentPhone    string
	ParentTg       string

	TgAdminUsername string
}
