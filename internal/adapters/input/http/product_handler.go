package http

import (
	"net/http"
	"strings"

	"project/internal/domain"
	"project/internal/ports/input"

	"github.com/gin-gonic/gin"
)

// ProductHandler maneja las peticiones HTTP para productos.
type ProductHandler struct {
	service input.ProductServicePort
}

// NewProductHandler crea una nueva instancia de ProductHandler.
func NewProductHandler(service input.ProductServicePort) *ProductHandler {
	return &ProductHandler{service: service}
}

// RegisterRoutes registra las rutas del handler en el router de Gin.
func (h *ProductHandler) RegisterRoutes(r *gin.Engine) {

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"version": "1.0.0",
		})
	})

	products := r.Group("/products")
	{
		products.GET("", h.GetAllProducts)
		products.GET("/compare", h.CompareProducts)
		products.GET("/:id", h.GetProductByID)
	}
}

// GetAllProducts godoc
// @Summary      Lista todos los productos
// @Description  Retorna una lista completa de productos disponibles
// @Tags         products
// @Produce      json
// @Success      200  {array}   domain.Product
// @Failure      500  {object}  ErrorResponse
// @Router       /products [get]
func (h *ProductHandler) GetAllProducts(c *gin.Context) {
	products, err := h.service.GetAllProducts()
	if err != nil {
		handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, products)
}

// GetProductByID godoc
// @Summary      Obtiene un producto por ID
// @Description  Retorna un producto específico dado su ID
// @Tags         products
// @Produce      json
// @Param        id   path      string  true  "Product ID"
// @Success      200  {object}  domain.Product
// @Failure      404  {object}  ErrorResponse
// @Failure      500  {object}  ErrorResponse
// @Router       /products/{id} [get]
func (h *ProductHandler) GetProductByID(c *gin.Context) {
	id := c.Param("id")
	product, err := h.service.GetProduct(id)
	if err != nil {
		handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, product)
}

// CompareProducts godoc
// @Summary      Compara múltiples productos
// @Description  Retorna lista de productos para comparar dado una lista de IDs
// @Tags         products
// @Produce      json
// @Param        ids  query     string  true  "IDs separados por coma (ej: 1,2,3)"
// @Success      200  {array}   domain.Product
// @Failure      400  {object}  ErrorResponse
// @Failure      500  {object}  ErrorResponse
// @Router       /products/compare [get]
func (h *ProductHandler) CompareProducts(c *gin.Context) {
	idsParam := c.Query("ids")
	if idsParam == "" {
		handleError(c, domain.ErrInvalidInput)
		return
	}

	ids := strings.Split(idsParam, ",")
	if len(ids) < 2 {
		handleError(c, domain.ErrInvalidInput)
		return
	}

	products, err := h.service.CompareProducts(ids)
	if err != nil {
		handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, products)
}
