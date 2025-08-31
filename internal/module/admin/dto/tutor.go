package dto

type Tutor struct {
	ID          int64
	FullName    string
	Phone       string
	Tg          string
	CostPerHour string
	SubjectID   int64
	AdminID     int64
}

type TutorFinance struct {
	Conversion int64
	Count      int64
	Amount     int64
}
