package dao

import (
	"github.com/it-chep/tutors.git/internal/module/admin/dto"
	"github.com/it-chep/tutors.git/pkg/xo"
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
	Count  int64 `db:"count" json:"count"`
	Amount int64 `db:"amount" json:"amount"`
}

func (sf StudentFinance) ToDomain() dto.StudentFinance {
	return dto.StudentFinance{
		Count:  sf.Count,
		Amount: sf.Amount,
	}
}
