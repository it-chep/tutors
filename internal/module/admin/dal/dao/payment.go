package dao

import (
	"github.com/it-chep/tutors.git/internal/config"
	"github.com/it-chep/tutors.git/internal/module/admin/dto"
	"github.com/samber/lo"
)

type Payment struct {
	ID   int64  `db:"id"`
	Bank string `db:"bank"`
}

func (p *Payment) ToDomain() dto.Payment {
	return dto.Payment{
		ID:   p.ID,
		Bank: config.Bank(p.Bank),
	}
}

type Payments []Payment

func (p Payments) ToDomain() []dto.Payment {
	return lo.Map(p, func(item Payment, _ int) dto.Payment {
		return item.ToDomain()
	})
}
