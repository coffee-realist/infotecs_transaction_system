package handlers

import (
	"encoding/json"
	"errors"
	"github.com/coffee-realist/infotecs_transaction_system/internal/services"
	"github.com/joomcode/errorx"
	"log"
	"net/http"
)

// HTTPError отправляет HTTP-ответ с указанным статусом и сообщением об ошибке в формате JSON.
//
// Параметры:
//   - w: http.ResponseWriter для записи ответа.
//   - message: текст ошибки, который будет отправлен в теле ответа.
//   - statusCode: HTTP-статус, соответствующий ошибке.
func HTTPError(w http.ResponseWriter, message string, statusCode int) {
	log.Printf("HTTP %d: %s", statusCode, message)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

// handleServiceError обрабатывает ошибку, возвращаемую сервисом, и преобразует её в HTTP-ответ.
//
// Ошибка анализируется с использованием библиотеки errorx и сопоставляется с типами:
//   - services.IsNotFoundErr → HTTP 404
//   - services.IsClientErr   → HTTP 400
//   - иначе                  → HTTP 500
//
// Параметры:
//   - w: http.ResponseWriter для записи ответа.
//   - err: ошибка, возвращённая из сервисного слоя.
func handleServiceError(w http.ResponseWriter, err error) {
	var errorxErr *errorx.Error
	if errors.As(err, &errorxErr) {
		switch {
		case services.IsNotFoundErr(errorxErr):
			HTTPError(w, errorxErr.Error(), http.StatusNotFound)
		case services.IsClientErr(errorxErr):
			HTTPError(w, errorxErr.Error(), http.StatusBadRequest)
		default:
			HTTPError(w, "Internal server error", http.StatusInternalServerError)
		}
	}
}
