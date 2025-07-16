package main

import (
	"context"
	"database/sql"
	"errors"
	"github.com/coffee-realist/infotecs_transaction_system/internal/api"
	"github.com/coffee-realist/infotecs_transaction_system/internal/services"
	"github.com/coffee-realist/infotecs_transaction_system/internal/storage"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// Функция main - точка входа в приложение
func main() {
	// Подключаемся к базе данных
	db, err := storage.Connect()
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}

	// Гарантируем закрытие
	defer func() {
		err := db.Close()
		if err != nil {
			log.Fatalf("Failed to close the database: %v", err)
		}
	}()

	// Применяем миграции
	if err := applyMigrations(db); err != nil {
		log.Fatalf("Failed to apply migrations: %v", err)
	}

	// Создаем сервисы
	service := services.NewService(db)

	// Генерируем кошельки
	if err := service.TransferService.GenerateWallets(); err != nil {
		log.Fatalf("Failed to generate wallets: %v", err)
	}

	// Настраиваем маршруты
	handler := api.NewHandler(service.TransferService)
	router := handler.InitRoutes()

	// Запускаем сервер
	server := new(api.Server)
	go func() {
		if err := server.Run("8080", router); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Server failed: %v", err)
		}
	}()
	log.Println("Server started at :8080")

	// Ожидаем сигнал для завершения
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// Завершаем работу
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.ShutDown(ctx); err != nil {
		log.Fatalf("Server shutdown failed: %v", err)
	}
	log.Println("Server exited properly")
}

// applyMigrations применяет миграцию базы данных, используя предоставленное подключение к базе данных SQL.
// Обеспечивает актуальность схемы базы данных. Возвращает ошибку в случае сбоя миграции.
func applyMigrations(db *sql.DB) error {
	// Настройка миграции для sqlite
	driver, err := sqlite.WithInstance(db, &sqlite.Config{})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations", // Путь к миграциям
		"sqlite", driver)
	if err != nil {
		return err
	}

	// Применение миграций
	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}

	log.Println("Migrations applied successfully")
	return nil
}
