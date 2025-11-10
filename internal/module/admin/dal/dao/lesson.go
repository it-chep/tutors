package dao

import (
	"context"
	"fmt"
	"time"

	"github.com/it-chep/tutors.git/internal/module/admin/dto"
	"github.com/it-chep/tutors.git/pkg/xo"
	"github.com/samber/lo"
)

type LessonDefaultDAO struct {
	xo.ConductedLesson
}

func (l LessonDefaultDAO) ToDomain() dto.Lesson {
	less := dto.Lesson{
		ID:        l.ID,
		TutorID:   l.TutorID,
		StudentID: l.StudentID,
		Duration:  time.Duration(l.DurationInMinutes) * time.Minute,
		IsTrial:   l.IsTrial.Bool,
	}
	if l.CreatedAt.Valid {
		less.Date = l.CreatedAt.Time
	}
	return less
}

type LessonDAO struct {
	xo.ConductedLesson
	FirstName  string  `db:"first_name" json:"first_name"`
	LastName   string  `db:"last_name" json:"last_name"`
	MiddleName string  `db:"middle_name" json:"middle_name"`
	TutorName  *string `db:"tutor_name" json:"tutor_name"`
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
		TutorFullName: lo.FromPtr(s.TutorName),
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
