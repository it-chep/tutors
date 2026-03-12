package receipts

type Receipt struct {
	ID                string `json:"id"`
	Amount            string `json:"amount"`
	Comment           string `json:"comment"`
	CreatedAt         string `json:"created_at"`
	HasReceipt        bool   `json:"has_receipt"`
	ReceiptFileName   string `json:"receipt_file_name"`
	ReceiptUploadedAt string `json:"receipt_uploaded_at"`
}

type Response struct {
	Receipts []Receipt `json:"receipts"`
}
