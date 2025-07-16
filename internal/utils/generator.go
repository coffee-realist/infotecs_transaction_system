package utils

import (
	"crypto/rand"
	"encoding/hex"
)

// GenerateRandomAddress генерирует криптографически безопасный случайный адрес.
//
// Адрес состоит из 32 байт (256 бит), что соответствует 64 символам в шестнадцатеричном представлении.
//
// Возвращает:
//   - строку длиной 64 символа, содержащую hex-представление случайных байтов;
//   - ошибку, если не удалось сгенерировать случайные байты.
func GenerateRandomAddress() (string, error) {
	const addressLength = 32 // 32 байта = 64 hex-символа

	randomBytes := make([]byte, addressLength)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(randomBytes), nil
}
