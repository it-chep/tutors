package dto

import (
	"context"

	pkgContext "github.com/it-chep/tutors.git/pkg/context"
)

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

func (r Role) FrontString() string {
	switch r {
	case SuperAdminRole:
		return "super_admin"
	case AdminRole:
		return "admin"
	case TutorRole:
		return "tutor"
	default:
		return "unknown"
	}
}

func IsTutorRole(ctx context.Context) bool {
	roleID, _ := pkgContext.GetUserRole(ctx)
	if roleID == int8(TutorRole) {
		return true
	}
	return false
}

func IsAdminRole(ctx context.Context) bool {
	roleID, _ := pkgContext.GetUserRole(ctx)
	if roleID == int8(AdminRole) {
		return true
	}
	return false
}

func IsSuperAdminRole(ctx context.Context) bool {
	roleID, _ := pkgContext.GetUserRole(ctx)
	if roleID == int8(SuperAdminRole) {
		return true
	}
	return false
}
