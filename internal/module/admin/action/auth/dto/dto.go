package register_dto

import (
	"errors"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	Email string `json:"email"`
	jwt.RegisteredClaims
}

type User struct {
	Email    string `json:"email" db:"email"`
	Password string `json:"password" db:"password"`
}

type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (r RegisterRequest) Validate() error {
	if r.Email == "" || r.Password == "" {
		return errors.New("invalid request")
	}
	return nil
}

type VerifyRequest struct {
	Email string `json:"email"`
	Code  string `json:"code"`
}

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type CodeRegister struct {
	Password string
	Code     string
}
