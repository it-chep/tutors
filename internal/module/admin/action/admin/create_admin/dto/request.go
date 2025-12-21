package dto

import "github.com/it-chep/tutors.git/internal/module/admin/dto"

type CreateRequest struct {
	FullName string
	Email    string
	Tg       string
	Phone    string

	Role         dto.Role
	AvailableTGs []string
}
