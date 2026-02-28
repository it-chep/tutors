package dao

import (
	"github.com/it-chep/tutors.git/internal/module/admin/dto"
	"github.com/samber/lo"
)

type TgAdminUsername struct {
	ID   int64  `db:"id"`
	Name string `db:"name"`
}

func (t TgAdminUsername) ToDomain() dto.TgAdminUsername {
	return dto.TgAdminUsername{
		ID:   t.ID,
		Name: t.Name,
	}
}

type TgAdminUsernames []TgAdminUsername

func (t TgAdminUsernames) ToDomain() dto.TgAdminUsernames {
	return lo.Map(t, func(item TgAdminUsername, _ int) dto.TgAdminUsername {
		return item.ToDomain()
	})
}
