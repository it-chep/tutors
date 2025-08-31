package dao

import (
	"github.com/it-chep/tutors.git/internal/module/admin/dto"
	"github.com/it-chep/tutors.git/pkg/xo"
)

type AdminDAO struct {
	xo.Admin
}

func (a AdminDAO) ToDomain() dto.Admin {
	return dto.Admin{
		ID:       a.ID,
		FullName: a.FullName,
		Phone:    a.Phone,
		Tg:       a.Tg,
	}
}
