package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/coffee-realist/infotecs_transaction_system/internal/services"
)

// GetTransactions обрабатывает HTTP-запрос для получения последних N транзакций.
//
// Возвращает http.HandlerFunc, который:
//   - Извлекает параметр "count" из строки запроса и преобразует его в целое число.
//   - Вызывает transferService.GetLastN для получения последних N транзакций.
//   - Возвращает JSON-массив транзакций при успехе или ошибку в формате JSON при сбое.
//
// Параметры:
//   - transferService: интерфейс, предоставляющий доступ к истории транзакций.
//
// Пример запроса:
//
//	GET 127.0.0.1:8080/api/transactions?count=2
func GetTransactions(transferService services.TransferInteractor) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Получаем параметр count
		countParam := r.URL.Query().Get("count")
		count, err := strconv.Atoi(countParam)
		if err != nil || count <= 0 {
			HTTPError(w, "Invalid count parameter", http.StatusBadRequest)
			return
		}

		// Получаем транзакции через сервис
		transactions, err := transferService.GetLastN(count)
		if err != nil {
			handleServiceError(w, err)
			return
		}

		// Формируем ответ
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(transactions); err != nil {
			HTTPError(w, "Failed to encode response", http.StatusInternalServerError)
		}
	}
}
