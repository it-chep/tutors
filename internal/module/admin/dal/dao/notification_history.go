package dao

import (
	"github.com/it-chep/tutors.git/internal/module/admin/dto"
	"github.com/it-chep/tutors.git/pkg/xo"
)

type NotificationHistoryDAO struct {
	xo.NotificationHistory
}

type NotificationsHistoryDAO []NotificationHistoryDAO

func (n *NotificationHistoryDAO) ToDomain() dto.NotificationHistory {
	return dto.NotificationHistory{
		ID:         n.ID,
		CreatedAt:  n.CreatedAt.Time,
		ParentTgID: n.ParentTgID.Int64,
		UserID:     n.UserID.Int64,
	}
}

func (n *NotificationsHistoryDAO) ToDomain() []dto.NotificationHistory {
	domain := make([]dto.NotificationHistory, 0, len(*n))
	for _, notification := range *n {
		domain = append(domain, notification.ToDomain())
	}
	return domain
}
