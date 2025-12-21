package dto

type User struct {
	ID int64

	Email    string
	Password string
	FullName string
	Tg       string
	Phone    string

	IsActive bool

	Role    Role
	AdminID int64
}

type UserInfo struct {
	ID      int64
	TutorID int64

	Role    Role
	AdminID int64
}

type PaidFunctions struct {
	AdminID       int64
	PaidFunctions map[string]bool
}
