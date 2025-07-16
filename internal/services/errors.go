package services

import (
	"github.com/joomcode/errorx"

	"github.com/coffee-realist/infotecs_transaction_system/internal/storage"
)

// IsNotFoundErr проверяет, является ли ошибка ошибкой "не найдено" (Not Found).
// Возвращает true, если ошибка или её причина соответствует ошибке NotFound из storage.
func IsNotFoundErr(err *errorx.Error) bool {
	return storage.IsNotFoundErr(err.Cause()) || storage.IsNotFoundErr(err)
}

// IsClientErr проверяет, является ли ошибка ошибкой клиента (Client error).
// Возвращает true, если ошибка содержит трейд Client или её причина - внешняя ошибка из storage.
func IsClientErr(err *errorx.Error) bool {
	return errorx.HasTrait(err, Client) || storage.IsExternalErr(err.Cause())
}

// IsServerErr проверяет, является ли ошибка серверной ошибкой (Server error).
// Возвращает true, если ошибка содержит трейд Server и её причина - внутренняя ошибка из storage.
func IsServerErr(err *errorx.Error) bool {
	return errorx.HasTrait(err, Server) && storage.IsInternalErr(err.Cause())
}

var (
	// ServiceErrors — пространство имён ошибок доменного слоя (сервисов).
	ServiceErrors = errorx.NewNamespace("domain")

	// Client — трейд для ошибок, связанных с некорректными действиями клиента.
	Client = errorx.RegisterTrait("client")
	// ErrInvalid — тип ошибки, указывающей на некорректные входные данные (ошибка клиента).
	ErrInvalid = ServiceErrors.NewType("invalid", Client)

	// Server — трейд для ошибок, связанных с внутренними ошибками сервера.
	Server = errorx.RegisterTrait("server")
	// ErrFailedToGenerate — тип ошибки, возникающей при ошибках генерации данных.
	ErrFailedToGenerate = ServiceErrors.NewType("failed_to_generate", Server)
	// ErrFailedToGet — тип ошибки при неудачном получении данных из хранилища.
	ErrFailedToGet = ServiceErrors.NewType("failed_to_get", Server)
	// ErrFailedToInsert — тип ошибки при ошибках вставки данных в хранилище.
	ErrFailedToInsert = ServiceErrors.NewType("failed_to_insert", Server)
	// ErrFailedToUpdate — тип ошибки при ошибках обновления данных (например, токена).
	ErrFailedToUpdate = ServiceErrors.NewType("failed_to_update token", Server)
)
