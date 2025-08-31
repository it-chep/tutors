package dto

type Student struct {
	ID         int64
	FirstName  string
	LastName   string
	MiddleName string
	Phone      string
	Tg         string

	CostPerHour string
	SubjectID   int64
	TutorID     int64

	ParentFullName string
	ParentPhone    string
	ParentTg       string

	IsFinishedTrial   bool
	IsBalanceNegative bool
	IsNewbie          bool
}

type StudentFinance struct {
	Count  int64
	Amount int64
}
