package payment_hash

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"io"
	"os"
	"strings"
)

// PaymentToken структура для хранения данных студента
type paymentToken struct {
	StudentID   int64  `json:"id"`
	StudentUUID string `json:"uuid"`
}

// DecryptPaymentHash расшифровывает хэш и возвращает данные студента
func DecryptPaymentHash(encryptedHash string) (studentID int64, studentUUID string, err error) {
	var token paymentToken
	secretKey := os.Getenv("PAYMENT_HASH_SECRET_KEY")

	// 1. Генерируем ключи из секрета
	encryptionKey, hmacKey := generateKeys(secretKey)

	// 2. Декодируем base64
	data, err := base64.URLEncoding.DecodeString(encryptedHash + "==")
	if err != nil {
		return token.StudentID, token.StudentUUID, fmt.Errorf("base64 decode error: %v", err)
	}

	if len(data) < 48 { // Минимальный размер (nonce + HMAC)
		return token.StudentID, token.StudentUUID, errors.New("invalid token length")
	}

	// 3. Разделяем на части
	// Формат: [12 байт nonce][шифртекст][32 байта HMAC]
	hmacSize := 32
	nonceSize := 12

	if len(data) < nonceSize+hmacSize {
		return token.StudentID, token.StudentUUID, errors.New("token too short")
	}

	// HMAC в конце
	hmacStart := len(data) - hmacSize
	ciphertextWithNonce := data[:hmacStart]
	receivedHMAC := data[hmacStart:]

	// 4. Проверяем HMAC
	h := hmac.New(sha256.New, hmacKey)
	h.Write(ciphertextWithNonce)
	expectedHMAC := h.Sum(nil)

	if !hmac.Equal(receivedHMAC, expectedHMAC) {
		return token.StudentID, token.StudentUUID, errors.New("invalid signature")
	}

	// 5. Разделяем nonce и шифртекст
	nonce := ciphertextWithNonce[:nonceSize]
	ciphertext := ciphertextWithNonce[nonceSize:]

	// 6. Создаем AES-GCM cipher
	block, err := aes.NewCipher(encryptionKey)
	if err != nil {
		return token.StudentID, token.StudentUUID, fmt.Errorf("cipher error: %v", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return token.StudentID, token.StudentUUID, fmt.Errorf("GCM error: %v", err)
	}

	// 7. Дешифруем
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return token.StudentID, token.StudentUUID, fmt.Errorf("decryption failed: %v", err)
	}

	// 8. Парсим JSON
	if err := json.Unmarshal(plaintext, &token); err != nil {
		return token.StudentID, token.StudentUUID, fmt.Errorf("invalid token data: %v", err)
	}

	return token.StudentID, token.StudentUUID, nil
}

// generateKeys создает два ключа из одного секрета
func generateKeys(secret string) ([]byte, []byte) {
	hash := sha256.New()

	// Ключ для шифрования
	hash.Write([]byte(secret + "-encryption-key-v1"))
	encryptionKey := make([]byte, 32)
	copy(encryptionKey, hash.Sum(nil))

	// Ключ для HMAC
	hash.Reset()
	hash.Write([]byte(secret + "-hmac-key-v1"))
	hmacKey := make([]byte, 32)
	copy(hmacKey, hash.Sum(nil))

	return encryptionKey, hmacKey
}

// EncryptPaymentData шифрует данные студента
func EncryptPaymentData(studentID int64, studentUUID string) (string, error) {
	secretKey := os.Getenv("PAYMENT_HASH_SECRET_KEY")

	encryptionKey, hmacKey := generateKeys(secretKey)

	token := paymentToken{
		StudentID:   studentID,
		StudentUUID: studentUUID,
	}

	plaintext, err := json.Marshal(token)
	if err != nil {
		return "", fmt.Errorf("marshal error: %v", err)
	}

	// 4. Создаем AES-GCM
	block, err := aes.NewCipher(encryptionKey)
	if err != nil {
		return "", fmt.Errorf("cipher error: %v", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("GCM error: %v", err)
	}

	// 5. Генерируем nonce
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", fmt.Errorf("nonce error: %v", err)
	}

	// 6. Шифруем
	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)

	// 7. Создаем HMAC
	h := hmac.New(sha256.New, hmacKey)
	h.Write(ciphertext)
	signature := h.Sum(nil)

	// 8. Объединяем: шифртекст + HMAC
	combined := append(ciphertext, signature...)

	// 9. Кодируем в base64
	encoded := base64.URLEncoding.EncodeToString(combined)

	// 10. Убираем padding для чистоты URL
	return strings.TrimRight(encoded, "="), nil
}
