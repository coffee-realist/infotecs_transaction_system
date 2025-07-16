package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/coffee-realist/infotecs_transaction_system/internal/dto"
	"github.com/coffee-realist/infotecs_transaction_system/internal/services"
)

// Send обрабатывает HTTP-запрос на перевод средств между кошельками.
//
// Возвращает http.HandlerFunc, который:
//   - Декодирует тело запроса в структуру dto.TransactionReq.
//   - Вызывает transferService.Send для выполнения перевода.
//   - Возвращает HTTP 201 (Created) при успехе или ошибку в формате JSON при сбое.
//
// Параметры:
//   - transferService: интерфейс, реализующий логику перевода средств.
//
// Пример успешного ответа:
//
//	HTTP/1.1 201 Created
//	Transaction successful
func Send(transferService services.TransferInteractor) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Декодируем запрос
		var req dto.TransactionReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			HTTPError(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		// Выполняем перевод через сервис
		if err := transferService.Send(req); err != nil {
			handleServiceError(w, err)
			return
		}

		// Формируем ответ
		w.WriteHeader(http.StatusCreated)
		_, err := w.Write([]byte("Transaction successful"))
		if err != nil {
			HTTPError(w, "Failed to write response", http.StatusInternalServerError)
		}
	}
}
