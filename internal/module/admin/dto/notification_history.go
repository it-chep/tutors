package dto

import (
	"time"
)

type NotificationHistory struct {
	ID         int64
	CreatedAt  time.Time
	ParentTgID int64
	UserID     int64
}
