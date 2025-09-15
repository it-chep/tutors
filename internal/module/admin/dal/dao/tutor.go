package dao

import (
	"database/sql"
	"github.com/it-chep/tutors.git/internal/module/admin/dto"
	"github.com/it-chep/tutors.git/internal/pkg/convert"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/samber/lo"
)

type TutorDAO struct {
	TutorID     sql.NullInt64 `db:"id"`
	CostPerHour string        `db:"cost_per_hour" json:"cost_per_hour"`
	SubjectID   int64         `db:"subject_id"`
	AdminID     int64         `db:"admin_id"`
	FullName    string        `db:"full_name"`
	Tg          string        `db:"tg"`
	Phone       string        `db:"phone"`
}

type TutorWithSubjectName struct {
	TutorDAO
	SubjectName sql.NullString `db:"subject_name"`
}

func (t TutorDAO) ToDomain() dto.Tutor {
	return dto.Tutor{
		ID:          t.TutorID.Int64,
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
	Conversion *int64          `db:"conversion" json:"conversion"`
	Count      *int64          `db:"count" json:"count"`
	Amount     *pgtype.Numeric `db:"amount" json:"amount"`
}

func (t TutorFinance) ToDomain() dto.TutorFinance {
	return dto.TutorFinance{
		Conversion: lo.FromPtr(t.Conversion),
		Count:      lo.FromPtr(t.Count),
		Amount:     convert.NumericToDecimal(lo.FromPtr(t.Amount)),
	}
}
