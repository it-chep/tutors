package dto

type Role int8

const (
	// UnknownRole - неизвестная роль
	UnknownRole Role = iota
	// SuperAdminRole - роль суперадмина
	SuperAdminRole
	// AdminRole - роль админа
	AdminRole
	// TutorRole - роль репетитора
	TutorRole
)

func (r Role) String() string {
	switch r {
	case SuperAdminRole:
		return "Супер админ"
	case AdminRole:
		return "Админ"
	case TutorRole:
		return "Репетитор"
	default:
		return "Неизвестно"
	}
}
