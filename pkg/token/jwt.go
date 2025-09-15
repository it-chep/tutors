package token

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	register_dto "github.com/it-chep/tutors.git/internal/module/admin/action/auth/dto"
)

func GenerateTokens(email, jwtKey, refreshKey string) (register_dto.TokenPair, error) {
	accessExp := time.Now().Add(15 * time.Minute)
	accessClaims := &register_dto.Claims{
		Email: email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(accessExp),
		},
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessString, err := accessToken.SignedString([]byte(jwtKey))
	if err != nil {
		return register_dto.TokenPair{}, err
	}

	refreshExp := time.Now().Add(14 * 24 * time.Hour)
	refreshClaims := &register_dto.Claims{
		Email: email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(refreshExp),
		},
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	_, err = refreshToken.SignedString([]byte(refreshKey))
	if err != nil {
		return register_dto.TokenPair{}, err
	}

	return register_dto.TokenPair{
		AccessToken: accessString,
		//RefreshToken: refreshString,
	}, nil
}
