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
		Tg:              TgLink(s.Tg),
		CostPerHour:     convert.NumericToDecimal(s.CostPerHour).String(),
		SubjectID:       s.SubjectID,
		TutorID:         s.TutorID,
		ParentFullName:  s.ParentFullName,
		ParentPhone:     s.ParentPhone,
		ParentTg:        TgLink(s.ParentTg),
		IsFinishedTrial: s.IsFinishedTrial,
		ParentTgID:      s.ParentTgID.Int64,
		TgAdminUsername: s.TgAdminUsername.String,
		IsArchived:      s.IsArchive.Bool,
		PaymentID:       s.PaymentID.Int64,
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

type ConductedLessonDAO struct {
	xo.ConductedLesson
}

type ConductedLessonDAOs []ConductedLessonDAO

type StudentTutorMoney struct {
	StudentID *int64          `db:"student_id"`
	TutorID   *int64          `db:"tutor_id"`
	Student   *pgtype.Numeric `db:"student_cost_per_hour" json:"student_cost_per_hour"`
	Tutor     *pgtype.Numeric `db:"tutor_cost_per_hour" json:"tutor_cost_per_hour"`
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

type Wallets []Wallet

func (w Wallet) ToDomain() dto.Wallet {
	return dto.Wallet{
		ID:        w.ID,
		StudentID: w.StudentID,
		Balance:   convert.NumericToDecimal(w.Balance),
	}
}

func (w Wallets) ToDomain() []dto.Wallet {
	return lo.Map(w, func(item Wallet, _ int) dto.Wallet {
		return item.ToDomain()
	})
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
