package add_available_tg

type Request struct {
	AvailableTg string `json:"tg_admin_username"`
	TgAdminID   *int64 `json:"tg_admin_id"`
}
