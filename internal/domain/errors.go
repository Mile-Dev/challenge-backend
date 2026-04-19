package domain

// AppError representa un error de dominio puro.
// No conoce nada de HTTP, gRPC ni ningún protocolo externo.
type AppError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// NewAppError crea un nuevo error de dominio.
func NewAppError(code, message string) *AppError {
	return &AppError{Code: code, Message: message}
}

// Error implementa la interfaz error.
func (e *AppError) Error() string {
	return e.Message
}

// Is permite comparar errores con errors.Is()
func (e *AppError) Is(target error) bool {
	t, ok := target.(*AppError)
	if !ok {
		return false
	}
	return e.Code == t.Code
}

// Errores de dominio tipados.
// Solo describen QUÉ salió mal en el negocio.
var (
	ErrProductNotFound = NewAppError(
		"product_not_found",
		"producto no encontrado",
	)

	ErrInvalidInput = NewAppError(
		"invalid_input",
		"parámetros de entrada inválidos",
	)

	ErrInternalServer = NewAppError(
		"internal_server_error",
		"error interno del servidor",
	)
)
