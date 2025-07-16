package dto

import "github.com/shopspring/decimal"

// BalanceResp представляет ответ с балансом кошелька.
type BalanceResp struct {
	Amount decimal.Decimal `json:"amount" db:"amount"` // Текущий баланс
}

// BalanceReq представляет запрос для получения баланса по адресу кошелька.
type BalanceReq struct {
	Address string `json:"address" db:"address"` // Адрес кошелька
}

// BalanceUpdateReq представляет запрос на обновление баланса кошелька.
type BalanceUpdateReq struct {
	Address string          `json:"address" db:"address"` // Адрес кошелька
	Amount  decimal.Decimal `json:"amount" db:"amount"`   // Новое значение баланса
}

// WalletReq представляет запрос с информацией о кошельке.
type WalletReq struct {
	Address string          `json:"address" db:"address"` // Адрес кошелька
	Balance decimal.Decimal `json:"balance" db:"balance"` // Баланс кошелька
}
