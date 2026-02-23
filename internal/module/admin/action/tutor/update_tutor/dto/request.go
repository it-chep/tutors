package dto

type UpdateRequest struct {
	FullName        string
	Phone           string
	Tg              string
	CostPerHour     string
	SubjectID       int64
	TgAdminUsername string
}
