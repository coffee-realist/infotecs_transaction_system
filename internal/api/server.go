package api

import (
	"context"
	"net/http"
	"time"
)

// Server инкапсулирует HTTP-сервер и методы управления его жизненным циклом.
type Server struct {
	httpServer *http.Server
}

// Run запускает HTTP-сервер на указанном порту с заданным обработчиком.
//
// Аргументы:
//   - port: строка с номером порта (например, "8080").
//   - handler: http.Handler, обрабатывающий входящие запросы.
//
// Возвращает:
//   - ошибку, если запуск сервера завершился с ошибкой.
func (s *Server) Run(port string, handler http.Handler) error {
	s.httpServer = &http.Server{
		Addr:           ":" + port,
		Handler:        handler,
		MaxHeaderBytes: 1 << 20,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
	}
	return s.httpServer.ListenAndServe()
}

// ShutDown завершает работу HTTP-сервера с учётом контекста завершения.
//
// Аргументы:
//   - ctx: контекст, управляющий таймаутом завершения сервера.
//
// Возвращает:
//   - ошибку при завершении, если такая возникла.
func (s *Server) ShutDown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
