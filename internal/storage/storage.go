package storage

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

// Connect устанавливает соединение с базой данных SQLite.
//
// Возвращает:
//   - указатель на объект базы данных *sql.DB при успешном подключении.
//   - ошибку, если не удалось открыть соединение с базой данных.
func Connect() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "./database.db")
	if err != nil {
		log.Printf("Error opening database: %v", err)
		return nil, err
	}
	return db, nil
}

// Repository агрегирует репозитории для работы с кошельками и транзакциями.
//
// Содержит интерфейсы WalletStorageInteractor и TransactionStorageInteractor,
// обеспечивающие доступ к методам хранения и извлечения данных.
type Repository struct {
	WalletRepository      WalletStorageInteractor
	TransactionRepository TransactionStorageInteractor
}

// NewRepository создаёт новый экземпляр Repository, инициализируя вложенные репозитории.
//
// Аргументы:
//   - db: подключение к базе данных *sql.DB.
//
// Возвращает:
//   - указатель на новый Repository, содержащий репозитории кошельков и транзакций.
func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		WalletRepository:      NewWalletRepository(db),
		TransactionRepository: NewTransactionRepository(db),
	}
}
