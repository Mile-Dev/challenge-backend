package output

import "project/internal/domain"

// ProductRepositoryPort define el contrato de salida para el repositorio de productos.
// Cualquier adaptador de datos debe implementar esta interfaz.
type ProductRepositoryPort interface {
	// FindByID busca un producto por su ID.
	// Retorna ErrProductNotFound si el producto no existe.
	FindByID(id string) (*domain.Product, error)

	// FindAll retorna todos los productos disponibles en el repositorio.
	FindAll() ([]*domain.Product, error)

	// FindByIDs retorna una lista de productos dado una lista de IDs.
	// Si algún ID no existe, retorna ErrProductNotFound.
	FindByIDs(ids []string) ([]*domain.Product, error)
}
