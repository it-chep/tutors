package dto

import "github.com/samber/lo"

type TgAdminUsername struct {
	ID   int64
	Name string
}

type TgAdminUsernames []TgAdminUsername

func (t TgAdminUsernames) IDs() []int64 {
	return lo.Map(t, func(item TgAdminUsername, _ int) int64 {
		return item.ID
	})
}
