package dto

import "time"

type UpdateLesson struct {
	Date     time.Time
	Duration time.Duration
}
