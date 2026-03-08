package dto

import "time"

type Comment struct {
	ID             int64
	UserID         int64
	StudentID      int64
	Text           string
	AuthorFullName string
	CreatedAt      time.Time
}
