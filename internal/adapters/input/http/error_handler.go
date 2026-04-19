package http

import (
	"errors"
	"log"
	"net/http"

	"project/internal/domain"

	"github.com/gin-gonic/gin"
)

// ErrorResponse representa la estructura de respuesta de error HTTP.
// Solo expone información segura al cliente.
type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// handleError traduce errores de dominio a respuestas HTTP seguras.
// Registra el error completo en logs internos pero solo
// expone información segura al cliente.
func handleError(c *gin.Context, err error) {
	// Log interno con contexto completo — solo visible en servidor
	log.Printf("[ERROR] %v", err)

	var appErr *domain.AppError
	if errors.As(err, &appErr) {
		// Al cliente solo llega el error de dominio limpio
		// Sin contexto interno ni detalles de implementación
		c.JSON(httpStatusFromError(appErr), ErrorResponse{
			Code:    appErr.Code,
			Message: appErr.Message,
		})
		return
	}

	// Error no esperado — nunca revelamos detalles internos
	c.JSON(http.StatusInternalServerError, ErrorResponse{
		Code:    domain.ErrInternalServer.Code,
		Message: domain.ErrInternalServer.Message,
	})
}

// httpStatusFromError mapea errores de dominio a códigos HTTP.
func httpStatusFromError(err *domain.AppError) int {
	switch err.Code {
	case domain.ErrProductNotFound.Code:
		return http.StatusNotFound
	case domain.ErrInvalidInput.Code:
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}
