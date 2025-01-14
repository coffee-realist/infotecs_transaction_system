package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

// GetLastTransactions извлекает самые последние транзакции из базы данных на основе указанного параметра count.
// Параметр count извлекается из строки запроса и определяет, сколько транзакций необходимо вернуть.
// Возвращает функцию обработчика HTTP, которая записывает транзакции в виде ответа в формате JSON или соответствующего сообщения об ошибке.
func GetLastTransactions(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Извлекаем параметр "count" из строки запроса
		countParam := r.URL.Query().Get("count")
		log.Println("Count parameter:", countParam)

		// Преобразуем строку в число. Если параметр некорректен, отправляем ошибку
		count, err := strconv.Atoi(countParam)
		if err != nil || count <= 0 {
			log.Println("Invalid count parameter:", countParam)
			http.Error(w, "Invalid count parameter", http.StatusBadRequest)
			return
		}

		log.Println("Fetching", count, "transactions")

		// Выполняем запрос для получения транзакций, сортируя по дате и ограничивая их количеством "count"
		rows, err := db.Query("SELECT from_address, to_address, amount, created_at FROM transactions ORDER BY created_at DESC LIMIT ?", count)
		if err != nil {
			log.Println("Failed to execute query:", err)
			http.Error(w, "Failed to retrieve transactions", http.StatusInternalServerError)
			return
		}
		// Обеспечиваем корректное закрытие объекта rows после использования
		defer func(rows *sql.Rows) {
			err := rows.Close()
			if err != nil {
				log.Println("Failed to close rows:", err)
			}
		}(rows)

		// Структура для хранения транзакций
		var transactions []struct {
			From      string  `json:"from"`
			To        string  `json:"to"`
			Amount    float64 `json:"amount"`
			CreatedAt string  `json:"created_at"`
		}
		// Итерация по строкам результата запроса и сохранение данных в срез
		for rows.Next() {
			var t struct {
				From      string  `json:"from"`
				To        string  `json:"to"`
				Amount    float64 `json:"amount"`
				CreatedAt string  `json:"created_at"`
			}
			// Сканируем данные из текущей строки в структуру
			if err := rows.Scan(&t.From, &t.To, &t.Amount, &t.CreatedAt); err != nil {
				log.Println("Failed to scan transaction:", err)
				http.Error(w, "Failed to scan transaction", http.StatusInternalServerError)
				return
			}
			transactions = append(transactions, t)
		}

		// Проверяем, не возникла ли ошибка при итерации по строкам
		if err := rows.Err(); err != nil {
			log.Println("Error iterating over rows:", err)
			http.Error(w, "Error iterating over rows", http.StatusInternalServerError)
			return
		}

		// Если транзакции не найдены, отправляем ошибку
		if len(transactions) == 0 {
			log.Println("No transactions found")
			http.Error(w, "No transactions found", http.StatusNotFound)
			return
		}

		// Устанавливаем заголовок Content-Type как JSON
		w.Header().Set("Content-Type", "application/json")
		// Кодируем список транзакций в формат JSON и отправляем клиенту
		if err := json.NewEncoder(w).Encode(transactions); err != nil {
			log.Println("Error encoding response:", err)
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		}
	}
}
