package jwt

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

var (
	jwtSecret    []byte
	jwtSecretKey = "JWT_SECRET_KEY"

	adminUsername    string
	adminUsernameKey = "ADMIN_NAME"

	adminPass    string
	adminPassKey = "ADMIN_PASSWORD"
)

func init() {
	if err := godotenv.Load(); err != nil {
		panic("env not set")
	}
	key := os.Getenv(jwtSecretKey)
	if key == "" {
		panic("JWT secret key not set")
	}
	jwtSecret = []byte(key)

	adminUsername = os.Getenv(adminUsernameKey)
	if adminUsername == "" {
		panic("admin user name not set")
	}

	adminPass = os.Getenv(adminPassKey)
	if adminPass == "" {
		panic("admin pass name not set")
	}
}

func GenerateJWT(username string) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)

	claims := &jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(expirationTime),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		Subject:   username,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func Valid(authToken string) bool {
	token, err := jwt.Parse(authToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return jwtSecret, nil
	})

	if err != nil || !token.Valid {
		return false
	}

	return true
}

// CheckCredentials проверяет пароль
func CheckCredentials(username, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(adminPass))
	return username == adminUsername && err == nil
}
