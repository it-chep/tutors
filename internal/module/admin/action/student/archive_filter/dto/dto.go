package dto

type FilterRequest struct {
	IsLost        bool
	TgUsernames   []string
	TgUsernameIDs []int64
	PaymentIDs    []int64
}
