package main

import (
	"database/sql"
	"infotecs_transactions_system/internal/database"
	"infotecs_transactions_system/internal/handlers"
	"infotecs_transactions_system/internal/services"
	"log"
	"net/http"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/gorilla/mux"
)

// main инициализирует подключение к базе данных, применяет миграции, генерирует кошельки, настраивает маршруты API и запускает сервер.
func main() {
	// Подключаемся к базе данных
	db, err := database.Connect()
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			log.Fatalf("Failed to close the database: %v", err)
		}
	}(db)

	// Применяем миграции
	if err := applyMigrations(db); err != nil {
		log.Fatalf("Failed to apply migrations: %v", err)
	}

	// Генерация кошельков при старте приложения, если они ещё не созданы
	if err := services.GenerateWallets(db); err != nil {
		log.Fatalf("Failed to generate wallets: %v", err)
	}

	// Настройка маршрутов
	router := mux.NewRouter()

	router.HandleFunc("/api/send", handlers.SendFunds(db)).Methods("POST")
	router.HandleFunc("/api/transactions", handlers.GetLastTransactions(db)).Methods("GET")
	router.HandleFunc("/api/wallet/{address}/balance", handlers.GetBalance(db)).Methods("GET")

	// Запуск сервера
	log.Println("Server started at :8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}

// applyMigrations применяет миграцию базы данных, используя предоставленное подключение к базе данных SQL.
// Обеспечивает актуальность схемы базы данных. Возвращает ошибку в случае сбоя миграции.
func applyMigrations(db *sql.DB) error {
	// Настройка миграции для sqlite
	driver, err := sqlite.WithInstance(db, &sqlite.Config{})
	if err != nil {
		return err
	}

	// Указываем путь к миграциям
	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations", // Путь к миграциям
		"sqlite", driver)
	if err != nil {
		return err
	}

	// Применение миграций
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}

	log.Println("Migrations applied successfully")
	return nil
}
