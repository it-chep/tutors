package dao

import (
	"github.com/it-chep/tutors.git/internal/module/admin/dto"
	"github.com/it-chep/tutors.git/internal/pkg/convert"
	"github.com/it-chep/tutors.git/pkg/xo"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/samber/lo"
)

type StudentDAO struct {
	xo.Student
}

func (s StudentDAO) ToDomain() dto.Student {
	return dto.Student{
		ID:              s.ID,
		FirstName:       s.FirstName,
		LastName:        s.LastName,
		MiddleName:      s.MiddleName,
		Phone:           s.Phone,
		Tg:              s.Tg,
		CostPerHour:     s.CostPerHour,
		SubjectID:       s.SubjectID,
		TutorID:         s.TutorID,
		ParentFullName:  s.ParentFullName,
		ParentPhone:     s.ParentPhone,
		ParentTg:        s.ParentTg,
		IsFinishedTrial: s.IsFinishedTrial,
		ParentTgID:      s.ParentTgID.Int64,
	}
}

type StudentsDAO []StudentDAO

func (studs StudentsDAO) ToDomain() []dto.Student {
	domain := make([]dto.Student, 0, len(studs))
	for _, student := range studs {
		domain = append(domain, student.ToDomain())
	}
	return domain
}

type StudentFinance struct {
	Count  *int64          `db:"count" json:"count"`
	Amount *pgtype.Numeric `db:"amount" json:"amount"`
}

func (sf StudentFinance) ToDomain() dto.StudentFinance {
	return dto.StudentFinance{
		Count:  lo.FromPtr(sf.Count),
		Amount: convert.NumericToDecimal(lo.FromPtr(sf.Amount)),
	}
}

type Wallet struct {
	xo.Wallet
}

func (w Wallet) ToDomain() dto.Wallet {
	return dto.Wallet{
		ID:        w.ID,
		StudentID: w.StudentID,
		Balance:   convert.NumericToDecimal(w.Balance),
	}
}

type StudentWithTransactions struct {
	StudentID         int64          `db:"student_id"`
	TutorID           int64          `db:"tutor_id"`
	IsFinishedTrial   bool           `db:"is_finished_trial"`
	TransactionsCount int64          `db:"transactions_count"`
	Balance           pgtype.Numeric `db:"balance"`
}

type StudentsWithTransactions []StudentWithTransactions

func (s StudentWithTransactions) ToDomain() dto.StudentWithTransactions {
	return dto.StudentWithTransactions{
		StudentID:         s.StudentID,
		TutorID:           s.TutorID,
		IsFinishedTrial:   s.IsFinishedTrial,
		TransactionsCount: s.TransactionsCount,
		Balance:           convert.NumericToDecimal(s.Balance),
	}
}

func (s StudentsWithTransactions) ToDomain() []dto.StudentWithTransactions {
	domain := make([]dto.StudentWithTransactions, 0, len(s))
	for _, student := range s {
		domain = append(domain, student.ToDomain())
	}
	return domain
}
