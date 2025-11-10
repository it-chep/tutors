package dto

import "time"

type Lesson struct {
	ID        int64
	Date      time.Time
	TutorID   int64
	StudentID int64
	Duration  time.Duration
	IsTrial   bool
}
