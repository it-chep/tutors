package dto

type User struct {
	ID int64

	Email    string
	Password string
	FullName string

	IsActive bool

	Role Role
}

type UserInfo struct {
	ID   int64
	Role Role
}
