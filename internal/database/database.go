package database

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

// Connect инициализирует подключение к базе данных SQLite и возвращает экземпляр базы данных или сообщение об ошибке в случае сбоя.
func Connect() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "./database.db")
	if err != nil {
		log.Printf("Error opening database: %v", err)
		return nil, err
	}
	return db, nil
}
