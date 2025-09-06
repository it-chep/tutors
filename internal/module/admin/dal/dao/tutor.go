package dao

import (
	"github.com/it-chep/tutors.git/internal/module/admin/dto"
	"github.com/it-chep/tutors.git/pkg/xo"
	"github.com/samber/lo"
)

type TutorDAO struct {
	xo.Tutor
}

func (t TutorDAO) ToDomain() dto.Tutor {
	return dto.Tutor{
		ID:          t.ID,
		FullName:    t.FullName,
		Phone:       t.Phone,
		Tg:          t.Tg,
		CostPerHour: t.CostPerHour,
		SubjectID:   t.SubjectID,
		AdminID:     t.AdminID,
	}
}

type TutorsDao []TutorDAO

func (ts TutorsDao) ToDomain() []dto.Tutor {
	domain := make([]dto.Tutor, 0, len(ts))
	for _, t := range ts {
		domain = append(domain, t.ToDomain())
	}
	return domain
}

type TutorFinance struct {
	Conversion *int64 `db:"conversion" json:"conversion"`
	Count      *int64 `db:"count" json:"count"`
	Amount     *int64 `db:"amount" json:"amount"`
}

func (t TutorFinance) ToDomain() dto.TutorFinance {
	return dto.TutorFinance{
		Conversion: lo.FromPtr(t.Conversion),
		Count:      lo.FromPtr(t.Count),
		Amount:     lo.FromPtr(t.Amount),
	}
}
