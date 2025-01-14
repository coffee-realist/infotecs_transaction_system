package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"net/http"
)

// GetBalance - это HTTP-обработчик, который извлекает баланс кошелька по адресу из базы данных и возвращает его в формате JSON.
// Адрес кошелька извлекается из URL-адреса запроса, и обработчик отправляет соответствующие коды состояния HTTP при обнаружении ошибок.
func GetBalance(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Извлекаем адрес кошелька из параметров URL
		address := mux.Vars(r)["address"]

		// Переменная для хранения баланса кошелька
		var balance float64
		// Выполняем запрос к базе данных для получения баланса
		err := db.QueryRow("SELECT balance FROM wallets WHERE address = ?", address).Scan(&balance)
		if err != nil {
			// Обработка ошибки, если кошелек не найден
			if errors.Is(err, sql.ErrNoRows) {
				// Если кошелек не найден в базе, отправляем ошибку 404
				http.Error(w, "Wallet not found", http.StatusNotFound)
			} else {
				// В случае других ошибок при запросе баланса отправляем 500 ошибку
				http.Error(w, "Failed to retrieve balance", http.StatusInternalServerError)
			}
			return
		}

		// Формируем ответ в формате JSON с балансом кошелька
		resp := map[string]float64{"balance": balance}
		// Устанавливаем заголовок Content-Type как JSON
		w.Header().Set("Content-Type", "application/json")
		// Кодируем ответ в JSON и отправляем клиенту
		err = json.NewEncoder(w).Encode(resp)
		if err != nil {
			// Если не удается закодировать ответ, ничего не делаем (можно добавить логирование ошибки)
			return
		}
	}
}
