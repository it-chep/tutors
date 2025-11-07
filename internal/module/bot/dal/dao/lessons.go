package dao

import (
	"github.com/it-chep/tutors.git/internal/module/bot/dto"
	"github.com/it-chep/tutors.git/pkg/xo"
	"time"
)

type LessonDAO struct {
	xo.ConductedLesson
}

type LessonsDAO []LessonDAO

func (l *LessonDAO) ToDomain() dto.Lesson {
	return dto.Lesson{
		ID:        l.ID,
		TutorID:   l.TutorID,
		StudentID: l.StudentID,
		Duration:  time.Duration(l.DurationInMinutes) * time.Minute,
		IsTrial:   l.IsTrial.Bool,
	}
}

func (lsd *LessonsDAO) ToDomain() []dto.Lesson {
	domain := make([]dto.Lesson, 0, len(*lsd))
	for _, l := range *lsd {
		domain = append(domain, l.ToDomain())
	}
	return domain
}
