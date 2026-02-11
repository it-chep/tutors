package payment_hash

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

// paymentToken структура для хранения данных
type paymentToken struct {
	StudentID   int64  `json:"sid"`
	StudentUUID string `json:"uuid"`
	Timestamp   int64  `json:"ts"`   // время создания
	Salt        string `json:"salt"` // случайная строка для уникальности
}

// DecryptPaymentHash расшифровывает хэш и возвращает данные студента
func DecryptPaymentHash(encryptedHash string) (studentID int64, studentUUID string, err error) {
	var token paymentToken

	// Получаем секретный ключ
	secretKey := os.Getenv("PAYMENT_HASH_SECRET_KEY")
	if secretKey == "" {
		return 0, "", errors.New("PAYMENT_HASH_SECRET_KEY is not set")
	}

	// 1. Генерируем ключи из секрета
	encryptionKey, hmacKey := generateKeys(secretKey)

	// 2. Декодируем base64 (добавляем padding если нужно)
	encryptedHash = addBase64Padding(encryptedHash)
	data, err := base64.URLEncoding.DecodeString(encryptedHash)
	if err != nil {
		return 0, "", fmt.Errorf("base64 decode error: %v", err)
	}

	// 3. Проверяем минимальный размер
	// AES-GCM nonce обычно 12 байт, HMAC-SHA256 32 байта, + минимум 1 байт данных
	if len(data) < 12+32+1 {
		return 0, "", errors.New("token too short")
	}

	// 4. Разделяем на части
	// Формат: [12 байт nonce][шифртекст][32 байта HMAC]
	hmacSize := 32
	hmacStart := len(data) - hmacSize

	ciphertextWithNonce := data[:hmacStart]
	receivedHMAC := data[hmacStart:]

	// 5. Проверяем HMAC
	h := hmac.New(sha256.New, hmacKey)
	h.Write(ciphertextWithNonce)
	expectedHMAC := h.Sum(nil)

	if !hmac.Equal(receivedHMAC, expectedHMAC) {
		return 0, "", errors.New("invalid signature")
	}

	// 6. Разделяем nonce и шифртекст
	// AES-GCM стандартный nonce size = 12 байт
	gcmNonceSize := 12
	if len(ciphertextWithNonce) < gcmNonceSize {
		return 0, "", errors.New("ciphertext too short for nonce")
	}

	nonce := ciphertextWithNonce[:gcmNonceSize]
	ciphertext := ciphertextWithNonce[gcmNonceSize:]

	// 7. Создаем AES-GCM cipher
	block, err := aes.NewCipher(encryptionKey)
	if err != nil {
		return 0, "", fmt.Errorf("cipher error: %v", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return 0, "", fmt.Errorf("GCM error: %v", err)
	}

	// 8. Дешифруем
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return 0, "", fmt.Errorf("decryption failed: %v", err)
	}

	// 9. Парсим JSON
	if err := json.Unmarshal(plaintext, &token); err != nil {
		return 0, "", fmt.Errorf("invalid token data: %v", err)
	}

	return token.StudentID, token.StudentUUID, nil
}

// EncryptPaymentData шифрует данные студента
func EncryptPaymentData(studentID int64, studentUUID string) (string, error) {
	// Получаем секретный ключ
	secretKey := os.Getenv("PAYMENT_HASH_SECRET_KEY")
	if secretKey == "" {
		return "", errors.New("PAYMENT_HASH_SECRET_KEY is not set")
	}

	encryptionKey, hmacKey := generateKeys(secretKey)

	// Создаем токен с timestamp
	token := paymentToken{
		StudentID:   studentID,
		StudentUUID: studentUUID,
		Timestamp:   time.Now().Unix(),
		Salt:        generateRandomString(8), // добавляем случайность
	}

	// Сериализуем в JSON
	plaintext, err := json.Marshal(token)
	if err != nil {
		return "", fmt.Errorf("marshal error: %v", err)
	}

	// Создаем AES-GCM
	block, err := aes.NewCipher(encryptionKey)
	if err != nil {
		return "", fmt.Errorf("cipher error: %v", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("GCM error: %v", err)
	}

	// Генерируем nonce
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", fmt.Errorf("nonce error: %v", err)
	}

	// Шифруем
	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)

	// Создаем HMAC подпись
	h := hmac.New(sha256.New, hmacKey)
	h.Write(ciphertext)
	signature := h.Sum(nil)

	// Объединяем: шифртекст + HMAC
	combined := append(ciphertext, signature...)

	// Кодируем в base64 URL-safe
	encoded := base64.URLEncoding.EncodeToString(combined)

	// Убираем padding для URL
	return strings.TrimRight(encoded, "="), nil
}

// generateKeys создает два ключа из одного секрета
func generateKeys(secret string) ([]byte, []byte) {
	// Используем разные контексты для генерации ключей
	encryptionKey := deriveKey(secret, "encryption-key")
	hmacKey := deriveKey(secret, "hmac-key")
	return encryptionKey, hmacKey
}

// deriveKey создает ключ фиксированной длины из секрета
func deriveKey(secret, context string) []byte {
	hash := sha256.New()
	hash.Write([]byte(secret + "-" + context + "-v1.0"))
	return hash.Sum(nil) // 32 байта
}

// generateRandomString создает случайную строку
func generateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)

	// Используем crypto/rand для криптографической безопасности
	if _, err := rand.Read(b); err != nil {
		// Fallback на time-based случайность
		nano := time.Now().UnixNano()
		for i := range b {
			b[i] = charset[int(nano>>uint(i*8))%len(charset)]
		}
	} else {
		for i := range b {
			b[i] = charset[int(b[i])%len(charset)]
		}
	}

	return string(b)
}

// addBase64Padding добавляет padding к base64 строке если нужно
func addBase64Padding(s string) string {
	// Base64 требует длина кратной 4
	switch len(s) % 4 {
	case 2:
		s += "=="
	case 3:
		s += "="
	}
	return s
}
