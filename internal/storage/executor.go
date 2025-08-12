package storage

import "database/sql"

// DBExecutor определяет интерфейс для выполнения SQL-запросов и подготовки операторов.
//
// Интерфейс абстрагирует методы, необходимые для выполнения команд, запросов и подготовки
// SQL-операторов, чтобы можно было работать как с *sql.DB, так и с транзакциями *sql.Tx.
type DBExecutor interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
	Prepare(query string) (*sql.Stmt, error)
}
