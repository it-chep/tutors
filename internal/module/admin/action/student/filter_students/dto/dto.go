package dto

type FilterRequest struct {
	IsLost        bool
	TgUsernameIDs []int64
	TgUsernames   []string
	PaymentIDs    []int64
}
