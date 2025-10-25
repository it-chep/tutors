package dao

import (
	"context"
	"fmt"
	"time"

	"github.com/it-chep/tutors.git/internal/module/admin/dto"
	"github.com/it-chep/tutors.git/pkg/xo"
	"github.com/samber/lo"
)

type LessonDAO struct {
	xo.ConductedLesson
	FirstName  string `db:"first_name" json:"first_name"`
	LastName   string `db:"last_name" json:"last_name"`
	MiddleName string `db:"middle_name" json:"middle_name"`
}

func (s *LessonDAO) ToDomain(ctx context.Context) dto.Lesson {
	l := dto.Lesson{
		ID:        s.ID,
		TutorID:   s.TutorID,
		StudentID: s.StudentID,
		Duration:  time.Duration(s.DurationInMinutes) * time.Minute,
		IsTrial:   s.IsTrial.Bool,
		StudentFullName: lo.Ternary(dto.IsTutorRole(ctx),
			fmt.Sprintf("%s %s", s.FirstName, s.MiddleName),
			fmt.Sprintf("%s %s %s", s.LastName, s.FirstName, s.MiddleName)),
	}
	if s.CreatedAt.Valid {
		l.Date = s.CreatedAt.Time
	}
	return l
}

type LessonsDAO []*LessonDAO

func (les LessonsDAO) ToDomain(ctx context.Context) []dto.Lesson {
	domain := make([]dto.Lesson, 0, len(les))
	for _, l := range les {
		domain = append(domain, l.ToDomain(ctx))
	}
	return domain
}
