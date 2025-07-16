package dto

import "time"
import "github.com/shopspring/decimal"

// TransactionReq представляет запрос на выполнение транзакции.
type TransactionReq struct {
	From   string          `json:"from" db:"from_address"` // Адрес отправителя
	To     string          `json:"to" db:"to_address"`     // Адрес получателя
	Amount decimal.Decimal `json:"amount" db:"amount"`     // Сумма перевода
}

// TransactionResp представляет ответ с информацией о транзакции.
type TransactionResp struct {
	From      string          `json:"from" db:"from_address"`     // Адрес отправителя
	To        string          `json:"to" db:"to_address"`         // Адрес получателя
	Amount    decimal.Decimal `json:"amount" db:"amount"`         // Сумма перевода
	CreatedAt time.Time       `json:"created_at" db:"created_at"` // Время создания транзакции
}

// TransactionsResp представляет список транзакций.
type TransactionsResp []TransactionResp
