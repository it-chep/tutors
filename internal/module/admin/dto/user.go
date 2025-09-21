package dto

type User struct {
	ID int64

	Email    string
	Password string
	FullName string
	Tg       string
	Phone    string

	IsActive bool

	Role Role
}

type UserInfo struct {
	ID      int64
	TutorID int64

	Role Role
}
