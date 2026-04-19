package input

import "project/internal/domain"

// ProductServicePort define el contrato de entrada para el servicio de productos.
// Cualquier adaptador HTTP debe usar este puerto para interactuar con la aplicación.
type ProductServicePort interface {
	// GetProduct retorna un producto por su ID.
	GetProduct(id string) (*domain.Product, error)

	// GetAllProducts retorna todos los productos disponibles.
	GetAllProducts() ([]*domain.Product, error)

	// CompareProducts retorna una lista de productos dado una lista de IDs.
	CompareProducts(ids []string) ([]*domain.Product, error)
}
