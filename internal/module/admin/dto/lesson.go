package dto

import "time"

type Lesson struct {
	ID int64

	Duration time.Duration
	IsTrial  bool

	Date      time.Time
	TutorID   int64
	StudentID int64

	StudentFullName string
	TutorFullName   string
}
