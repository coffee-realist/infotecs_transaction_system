package storage

import (
	"database/sql"
	_ "embed"
	"errors"
	"github.com/coffee-realist/infotecs_transaction_system/internal/dto"
)

// WalletsRepository реализует методы взаимодействия с таблицей кошельков в базе данных.
//
// Использует DBExecutor для выполнения SQL-запросов.
type WalletsRepository struct {
	executor DBExecutor
}

// NewWalletRepository создаёт новый экземпляр WalletsRepository.
//
// Аргументы:
//   - executor: объект, реализующий интерфейс DBExecutor.
//
// Возвращает:
//   - указатель на WalletsRepository.
func NewWalletRepository(executor DBExecutor) *WalletsRepository {
	return &WalletsRepository{executor: executor}
}

// WalletStorageInteractor описывает интерфейс работы с кошельками.
type WalletStorageInteractor interface {
	// Insert создаёт новый кошелёк.
	Insert(req dto.WalletReq) error
	// GetBalance возвращает баланс по адресу кошелька.
	GetBalance(req dto.BalanceReq) (dto.BalanceResp, error)
	// UpdateBalance обновляет баланс указанного кошелька.
	UpdateBalance(req dto.BalanceUpdateReq) error
	// GetCount возвращает количество зарегистрированных кошельков.
	GetCount() (int, error)
}

var (
	//go:embed assets/wallets/insert.sql
	walletsInsertSQL string

	//go:embed assets/wallets/get_balance.sql
	walletsGetBalanceSQL string

	//go:embed assets/wallets/get_count.sql
	walletsGetCountSQL string

	//go:embed assets/wallets/update_balance.sql
	walletsUpdateSQL string
)

// Insert добавляет новый кошелёк в базу данных.
//
// Аргументы:
//   - req: структура dto.WalletReq с адресом и начальным балансом.
//
// Возвращает:
//   - ошибку, если операция завершилась неуспешно.
func (r *WalletsRepository) Insert(req dto.WalletReq) error {
	_, err := r.executor.Exec(walletsInsertSQL, req.Address, req.Balance)
	if err != nil {
		return ErrFailedToInsert.Wrap(err, "failed to create wallet: %v", err)
	}
	return nil
}

// GetBalance возвращает текущий баланс кошелька по его адресу.
//
// Аргументы:
//   - req: структура dto.BalanceReq с адресом кошелька.
//
// Возвращает:
//   - dto.BalanceResp с текущим балансом.
//   - ошибку, если кошелёк не найден или возникла другая проблема.
func (r *WalletsRepository) GetBalance(req dto.BalanceReq) (dto.BalanceResp, error) {
	var balance dto.BalanceResp
	err := r.executor.QueryRow(walletsGetBalanceSQL, req.Address).Scan(&balance.Amount)
	if err != nil {
		// Обработка ошибки, если кошелек не найден или иная проблема
		if errors.Is(err, sql.ErrNoRows) {
			return dto.BalanceResp{}, ErrNotFound.Wrap(err, "wallet not found")
		} else {
			return dto.BalanceResp{}, ErrFailedToGet.Wrap(err, "failed to get balance")
		}
	}
	return balance, nil
}

// GetCount возвращает количество кошельков в системе.
//
// Возвращает:
//   - количество кошельков.
//   - ошибку при неудачном выполнении запроса.
func (r *WalletsRepository) GetCount() (int, error) {
	var count int
	err := r.executor.QueryRow(walletsGetCountSQL).Scan(&count)
	if err != nil {
		return -1, ErrFailedToGet.Wrap(err, "failed to get count of wallets")
	}

	return count, nil
}

// UpdateBalance обновляет баланс конкретного кошелька по адресу.
//
// Аргументы:
//   - req: структура dto.BalanceUpdateReq с адресом и изменением баланса.
//
// Возвращает:
//   - ошибку, если операция обновления завершилась неуспешно или кошелёк не найден.
func (r *WalletsRepository) UpdateBalance(req dto.BalanceUpdateReq) error {
	res, err := r.executor.Exec(walletsUpdateSQL, req.Amount, req.Address)
	if err != nil {
		return ErrFailedToUpdate.Wrap(err, "failed to update balance")
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return ErrFailedToUpdate.Wrap(err, "failed to get rows affected")
	}

	if rowsAffected == 0 {
		return ErrNotFound.New("wallet not found")
	}

	return nil
}
