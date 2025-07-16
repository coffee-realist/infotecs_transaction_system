package services

import (
	"database/sql"
	"github.com/coffee-realist/infotecs_transaction_system/internal/storage"
)

// Service агрегирует основные сервисы приложения.
type Service struct {
	TransferService TransferInteractor
}

// NewService создаёт и возвращает новый экземпляр Service,
// инициализируя все необходимые сервисы и репозитории.
// В качестве параметра принимает подключение к базе данных *sql.DB.
func NewService(db *sql.DB) *Service {
	repository := storage.NewRepository(db)
	return &Service{
		TransferService: NewTransferService(db, repository.WalletRepository, repository.TransactionRepository),
	}
}
