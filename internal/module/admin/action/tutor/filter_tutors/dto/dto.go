package dto

type FilterRequest struct {
	TgUsernameIDs []int64
	IsArchive     bool
}
