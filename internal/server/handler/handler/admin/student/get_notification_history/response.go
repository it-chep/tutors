package get_notification_history

type Notification struct {
	ID        int64  `json:"id"`
	CreatedAt string `json:"created_at"`
}

type Response struct {
	Notifications     []Notification `json:"notifications"`
	NotificationCount int64          `json:"notifications_count"`
}
