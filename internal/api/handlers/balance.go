package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/coffee-realist/infotecs_transaction_system/internal/dto"
	"github.com/coffee-realist/infotecs_transaction_system/internal/services"
)

// GetBalance возвращает http.HandlerFunc, обрабатывающий запрос на получение баланса кошелька по адресу.
//
// Параметры:
//   - transferService: интерфейс, предоставляющий метод получения баланса.
//
// Обрабатывает GET-запрос по пути, содержащему параметр "address".
// В случае успеха возвращает JSON-ответ с балансом, в противном случае — сообщение об ошибке.
//
// Пример пути: 127.0.0.1:8080/api/wallet/{address}/balance
func GetBalance(transferService services.TransferInteractor) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Извлекаем адрес из пути
		address := r.PathValue("address")
		if address == "" {
			HTTPError(w, "Wallet address is required", http.StatusBadRequest)
			return
		}

		// Получаем баланс через сервис
		balanceResp, err := transferService.GetBalance(dto.BalanceReq{Address: address})
		if err != nil {
			handleServiceError(w, err)
			return
		}

		// Формируем ответ
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(balanceResp); err != nil {
			HTTPError(w, "Failed to encode response", http.StatusInternalServerError)
		}
	}
}
