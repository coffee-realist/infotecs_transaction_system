package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
)

// SendRequest представляет собой запрос на перевод денежных средств между двумя юридическими лицами с указанием отправителя, получателя и суммы.
type SendRequest struct {
	From   string  `json:"from"`
	To     string  `json:"to"`
	Amount float64 `json:"amount"`
}

// SendFunds управляет переводом средств между кошельками. Он выполняет проверку, обновляет баланс и регистрирует транзакцию.
func SendFunds(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req SendRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		// Проверяем, чтобы все поля были корректны
		if req.From == "" || req.To == "" || req.Amount <= 0 {
			http.Error(w, "Invalid input fields", http.StatusBadRequest)
			return
		}

		// Проверяем, чтобы отправитель и получатель не были одним и тем же кошельком
		if req.From == req.To {
			http.Error(w, "Cannot transfer funds to the same wallet", http.StatusBadRequest)
			return
		}

		// Проверяем баланс отправителя **до** начала транзакции
		var balance float64
		err := db.QueryRow("SELECT balance FROM wallets WHERE address = ?", req.From).Scan(&balance)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				http.Error(w, "Sender wallet not found", http.StatusNotFound)
				log.Printf("Sender wallet %s not found", req.From)
			} else {
				http.Error(w, "Failed to retrieve sender's balance", http.StatusInternalServerError)
				log.Printf("Error querying sender balance: %v", err)
			}
			return
		}

		if balance < req.Amount {
			http.Error(w, "Insufficient funds", http.StatusBadRequest)
			log.Printf("Insufficient funds for wallet %s, balance: %.2f, attempted transfer: %.2f", req.From, balance, req.Amount)
			return
		}

		// Начинаем транзакцию только если все проверки пройдены
		tx, err := db.Begin()
		if err != nil {
			http.Error(w, "Failed to begin transaction", http.StatusInternalServerError)
			log.Printf("Error beginning transaction: %v", err) // Логирование ошибки начала транзакции
			return
		}

		// Отложенный вызов Rollback на случай, если транзакция не будет зафиксирована
		defer func() {
			if r := recover(); r != nil {
				_ = tx.Rollback() // Восстановление после паники
				panic(r)          // Перепаниковка
			} else if err != nil {
				_ = tx.Rollback()                                           // Откат транзакции в случае ошибки
				log.Printf("Transaction rolled back due to error: %v", err) // Логирование отката транзакции
			}
		}()

		// Обновляем баланс отправителя
		_, err = tx.Exec("UPDATE wallets SET balance = balance - ? WHERE address = ?", req.Amount, req.From)
		if err != nil {
			http.Error(w, "Failed to update sender's balance", http.StatusInternalServerError)
			log.Printf("Error updating sender's balance: %v", err)
			return
		}

		// Обновляем баланс получателя
		_, err = tx.Exec("UPDATE wallets SET balance = balance + ? WHERE address = ?", req.Amount, req.To)
		if err != nil {
			http.Error(w, "Failed to update receiver's balance", http.StatusInternalServerError)
			log.Printf("Error updating receiver's balance: %v", err)
			return
		}

		// Записываем транзакцию
		_, err = tx.Exec("INSERT INTO transactions (from_address, to_address, amount) VALUES (?, ?, ?)", req.From, req.To, req.Amount)
		if err != nil {
			http.Error(w, "Failed to record transaction", http.StatusInternalServerError)
			log.Printf("Error recording transaction: %v", err)
			return
		}

		// Фиксируем транзакцию
		err = tx.Commit()
		if err != nil {
			http.Error(w, "Failed to commit transaction", http.StatusInternalServerError)
			log.Printf("Error committing transaction: %v", err) // Логирование ошибки коммита
			return
		}

		// Успешный ответ
		w.WriteHeader(http.StatusCreated)
		_, err = w.Write([]byte("Transaction successful"))
		if err != nil {
			log.Printf("Failed to write response: %v", err)
			return
		}
		log.Printf("Transaction from %s to %s for amount %.2f\n", req.From, req.To, req.Amount)
	}
}
