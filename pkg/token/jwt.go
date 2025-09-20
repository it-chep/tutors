package token

import (
	"fmt"
	"net/http"
	"strings"
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
	refreshString, err := refreshToken.SignedString([]byte(refreshKey))
	if err != nil {
		return register_dto.TokenPair{}, err
	}

	return register_dto.NewTokenPair(accessString, refreshString), nil
}

func RefreshClaimsFromRequest(r *http.Request, refreshSecret string) (*register_dto.Claims, error) {
	cookie, err := r.Cookie(register_dto.RefreshCookie)
	if err != nil {
		return nil, err
	}

	claims := &register_dto.Claims{}
	token, err := jwt.ParseWithClaims(cookie.Value, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return []byte(refreshSecret), nil
	})

	if err != nil || !token.Valid || claims.ExpiresAt.Time.Before(time.Now()) {
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
}

func AccessClaimsFromRequest(r *http.Request, jwtAccessSecret string) (*register_dto.Claims, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return nil, fmt.Errorf("invalid token")
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || parts[0] != "Bearer" {
		return nil, fmt.Errorf("invalid token")
	}
	tokenStr := parts[1]

	claims := &register_dto.Claims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return []byte(jwtAccessSecret), nil
	})
	if err != nil || !token.Valid || claims.ExpiresAt.Time.Before(time.Now()) {
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
}
