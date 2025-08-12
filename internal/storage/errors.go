package storage

import "github.com/joomcode/errorx"

// IsInternalErr проверяет, является ли ошибка внутренней ошибкой хранилища.
//
// Аргументы:
//   - err: ошибка для проверки.
//
// Возвращает:
//   - true, если ошибка имеет признак Internal, иначе false.
func IsInternalErr(err error) bool {
	return errorx.HasTrait(err, Internal)
}

// IsExternalErr проверяет, является ли ошибка внешней ошибкой хранилища.
//
// Аргументы:
//   - err: ошибка для проверки.
//
// Возвращает:
//   - true, если ошибка имеет признак External, иначе false.
func IsExternalErr(err error) bool {
	return errorx.HasTrait(err, External)
}

// IsNotFoundErr проверяет, является ли ошибка ошибкой "не найдено".
//
// Аргументы:
//   - err: ошибка для проверки.
//
// Возвращает:
//   - true, если ошибка имеет признак NotFound, иначе false.
func IsNotFoundErr(err error) bool {
	return errorx.HasTrait(err, errorx.NotFound())
}

var (
	// Namespace пространство имён для ошибок слоя хранения.
	Namespace = errorx.NewNamespace("storage")

	// External признак внешних ошибок, связанных с хранилищем.
	External = errorx.RegisterTrait("external")
	// ErrNotFound ошибка "запись не найдена" с признаками External и NotFound.
	ErrNotFound = Namespace.NewType("not_found", External, errorx.NotFound())

	// Internal признак внутренних ошибок, связанных с хранилищем.
	Internal = errorx.RegisterTrait("internal")
	// ErrFailedToInsert ошибка при неудачной вставке данных.
	ErrFailedToInsert = Namespace.NewType("failed_to_insert", Internal)
	// ErrFailedToUpdate ошибка при неудачном обновлении данных.
	ErrFailedToUpdate = Namespace.NewType("failed_to_update", Internal)
	// ErrFailedToGet ошибка при неудачном получении данных.
	ErrFailedToGet = Namespace.NewType("failed_to_get", Internal)
	// ErrFailedToMarshal ошибка при ошибке маршалинга данных.
	ErrFailedToMarshal = Namespace.NewType("failed_to_marshal", Internal)
	// ErrFailedToUnmarshal ошибка при ошибке анмаршалинга данных.
	ErrFailedToUnmarshal = Namespace.NewType("failed_to_unmarshal", Internal)
	// ErrFailedToCloseRows ошибка при неудачном закрытии строк результата запроса.
	ErrFailedToCloseRows = Namespace.NewType("failed_to_close_rows", Internal)
	// UnhandledErr ошибка, не обработанная явно.
	UnhandledErr = Namespace.NewType("unhandled", Internal)
)
