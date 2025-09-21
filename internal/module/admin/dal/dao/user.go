package dao

import (
	"github.com/it-chep/tutors.git/internal/module/admin/dto"
	"github.com/it-chep/tutors.git/pkg/xo"
)

type User struct {
	xo.User
}

func (u *User) UserDto() dto.User {
	return dto.User{
		ID:       u.ID,
		Email:    u.Email,
		Password: u.Password.String,
		Phone:    u.Phone,
		Tg:       TgLink(u.Tg),
		FullName: u.FullName.String,
		IsActive: u.IsActive.Bool,
		Role:     dto.Role(u.RoleID.Int64),
	}
}

func (u *User) UserInfo() *dto.UserInfo {
	return &dto.UserInfo{
		ID:      u.ID,
		TutorID: u.TutorID.Int64,
		Role:    dto.Role(u.RoleID.Int64),
	}
}

type Users []User

func (u Users) ToDomain() []dto.User {
	domain := make([]dto.User, 0, len(u))
	for _, user := range u {
		domain = append(domain, user.UserDto())
	}
	return domain
}
