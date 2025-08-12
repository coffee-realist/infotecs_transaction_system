package api

import (
	"github.com/coffee-realist/infotecs_transaction_system/internal/api/handlers"
	"github.com/coffee-realist/infotecs_transaction_system/internal/services"
	"net/http"
)

// Handler агрегирует зависимости для HTTP-обработчиков API.
type Handler struct {
	transferService services.TransferInteractor
}

// NewHandler создаёт новый экземпляр Handler.
//
// Аргументы:
//   - transferService: реализация интерфейса TransferInteractor,
//     обеспечивающая доступ к операциям с транзакциями и кошельками.
//
// Возвращает:
//   - Указатель на Handler, содержащий переданный сервис.
func NewHandler(transferService services.TransferInteractor) *Handler {
	return &Handler{transferService: transferService}
}

// InitRoutes инициализирует маршруты HTTP-сервера.
//
// Регистрирует следующие маршруты:
//   - POST /api/send — отправка транзакции
//   - GET /api/transactions — получение последних N транзакций
//   - GET /api/wallet/{address}/balance — получение баланса кошелька
//
// Возвращает:
//   - http.Handler с зарегистрированными маршрутами.
func (h *Handler) InitRoutes() http.Handler {
	mux := http.NewServeMux()

	// Регистрируем обработчики с новым синтаксисом
	mux.HandleFunc("POST /api/send", handlers.Send(h.transferService))
	mux.HandleFunc("GET /api/transactions", handlers.GetTransactions(h.transferService))
	mux.HandleFunc("GET /api/wallet/{address}/balance", handlers.GetBalance(h.transferService))

	return mux
}
