package services

import (
	"database/sql"
	"github.com/coffee-realist/infotecs_transaction_system/internal/dto"
	"github.com/coffee-realist/infotecs_transaction_system/internal/storage"
	"github.com/coffee-realist/infotecs_transaction_system/internal/utils"
	"github.com/shopspring/decimal"
	"log"
)

// TransferInteractor описывает интерфейс бизнес-логики для операций с транзакциями и балансами.
type TransferInteractor interface {
	// Send выполняет перевод между кошельками согласно данным транзакции.
	Send(transaction dto.TransactionReq) error

	// GetLastN возвращает последние N транзакций.
	GetLastN(n int) (dto.TransactionsResp, error)

	// GetBalance возвращает баланс кошелька по адресу.
	GetBalance(req dto.BalanceReq) (dto.BalanceResp, error)

	// GenerateWallets инициализирует 10 кошельков с дефолтным балансом,
	// если в базе еще нет кошельков.
	GenerateWallets() error
}

// TransferService реализует TransferInteractor, используя репозитории и базу данных.
type TransferService struct {
	db                    *sql.DB
	walletRepository      storage.WalletStorageInteractor
	transactionRepository storage.TransactionStorageInteractor
}

// NewTransferService создаёт новый экземпляр TransferService.
//
// Аргументы:
//   - db: подключение к базе данных.
//   - walletRepository: репозиторий для работы с кошельками.
//   - transactionRepository: репозиторий для работы с транзакциями.
//
// Возвращает:
//   - Указатель на TransferService с инициализированными репозиториями и подключением к БД.
func NewTransferService(
	db *sql.DB,
	walletRepository storage.WalletStorageInteractor,
	transactionRepository storage.TransactionStorageInteractor,
) *TransferService {
	return &TransferService{
		db:                    db,
		walletRepository:      walletRepository,
		transactionRepository: transactionRepository,
	}
}

// Send выполняет перевод средств между кошельками.
//
// Аргументы:
//   - req: структура dto.TransactionReq с полями From, To и Amount.
//
// Возвращает:
//   - Ошибку в случае некорректных данных, недостаточного баланса, ошибок работы с БД или репозиториями.
//     В случае успеха — nil.
func (s *TransferService) Send(req dto.TransactionReq) error {
	// Валидация
	if req.From == "" || req.To == "" || !req.Amount.IsPositive() {
		return ErrInvalid.New("invalid input fields")
	}
	if req.From == req.To {
		return ErrInvalid.New("cannot transfer to same wallet")
	}

	tx, err := s.db.Begin()
	if err != nil {
		return err
	}

	// Гарантируем откат при ошибках
	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				log.Printf("Rollback error: %v", rbErr)
			}
		}
	}()

	// Создаем репозитории, использующие транзакцию
	walletRepo := storage.NewWalletRepository(tx)
	transactionRepo := storage.NewTransactionRepository(tx)

	// Проверяем и обновляем баланс отправителя
	balanceResp, err := s.walletRepository.GetBalance(dto.BalanceReq{Address: req.From})
	if err != nil {
		return ErrFailedToGet.Wrap(err, "failed to get balance")
	}
	if balanceResp.Amount.LessThan(req.Amount) {
		return ErrInvalid.New("insufficient funds")
	}

	if err := walletRepo.UpdateBalance(dto.BalanceUpdateReq{
		Address: req.From,
		Amount:  balanceResp.Amount.Sub(req.Amount),
	}); err != nil {
		return ErrFailedToUpdate.Wrap(err, "failed to update balance")
	}

	// Проверяем и обновляем баланс получателя
	balanceResp, err = s.walletRepository.GetBalance(dto.BalanceReq{Address: req.To})
	if err != nil {
		return ErrFailedToGet.Wrap(err, "failed to get balance")
	}

	if err := walletRepo.UpdateBalance(dto.BalanceUpdateReq{
		Address: req.To,
		Amount:  balanceResp.Amount.Add(req.Amount),
	}); err != nil {
		return ErrFailedToUpdate.Wrap(err, "failed to update balance")
	}

	// Создаем запись о транзакции
	if err := transactionRepo.Insert(req); err != nil {
		return ErrFailedToInsert.Wrap(err, "failed to insert transaction")
	}

	// Фиксируем транзакцию
	if err := tx.Commit(); err != nil {
		return ErrFailedToInsert.Wrap(err, "failed to commit transaction")
	}

	return nil
}

// GenerateWallets создаёт 10 кошельков с начальным балансом 100, если в базе нет кошельков.
//
// Возвращает:
//   - Ошибку при сбое проверки существующих кошельков, генерации адресов или записи в базу.
//     В случае успеха — nil.
func (s *TransferService) GenerateWallets() error {
	// Проверяем наличие уже существующих кошельков
	count, err := s.walletRepository.GetCount()
	if err != nil {
		return ErrFailedToGet.WrapWithNoMessage(err)
	}
	if count > 0 {
		return nil // Кошельки уже существуют
	}
	// Генерируем 10 кошельков
	for i := 0; i < 10; i++ {
		address, err := utils.GenerateRandomAddress()
		if err != nil {
			return ErrFailedToGenerate.Wrap(err, "failed to generate wallet address")
		}
		if err = s.walletRepository.Insert(dto.WalletReq{Address: address, Balance: decimal.NewFromInt(100)}); err != nil {
			return ErrFailedToInsert.WrapWithNoMessage(err)
		}
	}

	return nil
}

// GetLastN возвращает последние N транзакций.
//
// Аргументы:
//   - n: количество транзакций для выборки.
//
// Возвращает:
//   - Срез транзакций dto.TransactionsResp и ошибку, если она возникла.
func (s *TransferService) GetLastN(n int) (dto.TransactionsResp, error) {
	return s.transactionRepository.GetLastN(n)
}

// GetBalance возвращает баланс кошелька по адресу.
//
// Аргументы:
//   - req: структура dto.BalanceReq с адресом кошелька.
//
// Возвращает:
//   - Баланс dto.BalanceResp и ошибку при её возникновении.
func (s *TransferService) GetBalance(req dto.BalanceReq) (dto.BalanceResp, error) {
	return s.walletRepository.GetBalance(req)
}
