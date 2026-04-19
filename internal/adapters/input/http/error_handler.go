package http

import (
	"errors"
	"log/slog"
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

	var appErr *domain.AppError
	if errors.As(err, &appErr) {

		// Logueamos el error completo con contexto para diagnóstico interno
		slog.Error("request error",
			"code", appErr.Code,
			"message", appErr.Message,
			"path", c.Request.URL.Path,
			"method", c.Request.Method,
		)
		// Al cliente solo llega el error de dominio limpio
		// Sin contexto interno ni detalles de implementación
		c.JSON(httpStatusFromError(appErr), ErrorResponse{
			Code:    appErr.Code,
			Message: appErr.Message,
		})
		return
	}

	// Error inesperado — log con stack completo
	slog.Error("unexpected error",
		"error", err.Error(),
		"path", c.Request.URL.Path,
		"method", c.Request.Method,
	)
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
