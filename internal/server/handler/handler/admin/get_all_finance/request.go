package get_all_finance

type Request struct {
	From    string `json:"from"`
	To      string `json:"to"`
	AdminID int64  `json:"admin_id"`
}
