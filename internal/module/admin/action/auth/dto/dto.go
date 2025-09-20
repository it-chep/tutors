package register_dto

import (
	"errors"

	"github.com/golang-jwt/jwt/v5"
)

const (
	RefreshCookie = "100_rep_refresh"
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
	AccessToken  string `json:"token"`
	refreshToken string
}

func NewTokenPair(accessToken, refreshToken string) TokenPair {
	return TokenPair{
		AccessToken:  accessToken,
		refreshToken: refreshToken,
	}
}

func (p TokenPair) Refresh() string {
	return p.refreshToken
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type CodeRegister struct {
	Password string
	Code     string
}
