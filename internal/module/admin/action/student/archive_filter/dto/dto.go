package dto

type FilterRequest struct {
	IsLost      bool
	TgUsernames []string
}
