package storage

import (
	_ "embed"
	"github.com/coffee-realist/infotecs_transaction_system/internal/dto"
	"log"
)

// TransactionRepository реализует методы для работы с транзакциями в базе данных.
//
// Использует DBExecutor для выполнения SQL-запросов.
type TransactionRepository struct {
	executor DBExecutor
}

// NewTransactionRepository создаёт новый экземпляр TransactionRepository.
//
// Аргументы:
//   - executor: объект, реализующий интерфейс DBExecutor для выполнения запросов.
//
// Возвращает:
//   - указатель на TransactionRepository.
func NewTransactionRepository(executor DBExecutor) *TransactionRepository {
	return &TransactionRepository{executor: executor}
}

// TransactionStorageInteractor описывает интерфейс операций с транзакциями.
type TransactionStorageInteractor interface {
	// GetLastN возвращает последние n транзакций.
	GetLastN(n int) (dto.TransactionsResp, error)
	// Insert вставляет новую транзакцию в базу данных.
	Insert(req dto.TransactionReq) error
}

var (
	//go:embed assets/transactions/insert.sql
	transactionsInsertSQL string

	//go:embed assets/transactions/get_last_n.sql
	transactionsGetLastNSQL string
)

// GetLastN возвращает последние n транзакций из базы данных.
//
// Аргументы:
//   - n: количество транзакций для получения.
//
// Возвращает:
//   - срез транзакций dto.TransactionsResp при успешном выполнении.
//   - ошибку, если произошла проблема с запросом или данными.
func (r *TransactionRepository) GetLastN(n int) (dto.TransactionsResp, error) {
	var transactions dto.TransactionsResp

	// Получаем из БД последние n транзакций
	rows, err := r.executor.Query(transactionsGetLastNSQL, n)
	if err != nil {
		return dto.TransactionsResp{}, ErrFailedToGet.Wrap(err,
			"Failed to get last %d transactions", n)
	}

	// Обеспечиваем корректное закрытие объекта rows после использования
	defer func() {
		err := rows.Close()
		if err != nil {
			log.Println("Failed to close rows:", err)
		}
	}()

	// Итерация по строкам результата запроса и сохранение данных в срез
	for rows.Next() {
		var transaction dto.TransactionResp
		if err := rows.Scan(
			&transaction.From,
			&transaction.To,
			&transaction.Amount,
			&transaction.CreatedAt,
		); err != nil {
			return transactions, ErrFailedToUnmarshal.Wrap(err,
				"Failed to unmarshall last %d transactions", n)
		}
		transactions = append(transactions, transaction)
	}

	// Проверяем, не возникла ли ошибка при итерации по строкам
	if err := rows.Err(); err != nil {
		return dto.TransactionsResp{}, UnhandledErr.Wrap(err,
			"Error while scanning transactions through sql rows", n)
	}

	// Если транзакции не найдены, отправляем ошибку
	if len(transactions) == 0 {
		return dto.TransactionsResp{}, ErrNotFound.New("transactions not found")
	}

	return transactions, nil
}

// Insert добавляет новую транзакцию в базу данных.
//
// Аргументы:
//   - req: структура dto.TransactionReq с данными транзакции.
//
// Возвращает:
//   - ошибку, если не удалось вставить транзакцию в базу.
func (r *TransactionRepository) Insert(req dto.TransactionReq) error {
	_, err := r.executor.Exec(transactionsInsertSQL, req.From, req.To, req.Amount)
	if err != nil {
		return ErrFailedToInsert.Wrap(err, "failed to insert transaction: %v", err)
	}

	return nil
}
