package services

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"fmt"
)

// GenerateWallets инициализирует кошельки в базе данных, если ни один из них не существует, создавая 10 кошельков с балансом по умолчанию.
// Программа генерирует случайные адреса кошельков, вставляет их в таблицу "wallets" и возвращает ошибку, если какой-либо шаг не выполняется.
func GenerateWallets(db *sql.DB) error {
	// Проверяем наличие уже существующих кошельков
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM wallets").Scan(&count)
	if err != nil {
		return err
	}
	if count > 0 {
		return nil // Кошельки уже существуют
	}

	// Генерируем 10 кошельков
	for i := 0; i < 10; i++ {
		address, err := generateRandomAddress()
		if err != nil {
			return fmt.Errorf("failed to generate wallet address: %v", err)
		}
		_, err = db.Exec("INSERT INTO wallets (address, balance) VALUES (?, ?)", address, 100.0)
		if err != nil {
			return fmt.Errorf("failed to insert wallet into database: %v", err)
		}
	}
	return nil
}

// generateRandomAddress генерирует случайный 32-байтовый адрес и возвращает его в виде шестнадцатеричной строки. В случае неудачи возвращает сообщение об ошибке.
func generateRandomAddress() (string, error) {
	const addressLength = 32 // 32 байта = 64 символа в hex
	randomBytes := make([]byte, addressLength)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(randomBytes), nil
}
