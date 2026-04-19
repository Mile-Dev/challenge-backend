package json

import (
	"encoding/json"
	"os"

	"project/internal/domain"
	"project/internal/ports/output"
)

// ProductJSONRepository implementa ProductRepositoryPort usando un archivo JSON.
type ProductJSONRepository struct {
	products []*domain.Product
}

// NewProductJSONRepository crea una nueva instancia cargando productos desde un archivo JSON.
func NewProductJSONRepository(filePath string) (output.ProductRepositoryPort, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, domain.ErrInternalServer
	}

	var products []*domain.Product
	if err := json.Unmarshal(data, &products); err != nil {
		return nil, domain.ErrInternalServer
	}

	return &ProductJSONRepository{products: products}, nil
}

// FindByID busca un producto por su ID.
func (r *ProductJSONRepository) FindByID(id string) (*domain.Product, error) {
	for _, p := range r.products {
		if p.ID == id {
			return p, nil
		}
	}
	return nil, domain.ErrProductNotFound
}

// FindAll retorna todos los productos.
func (r *ProductJSONRepository) FindAll() ([]*domain.Product, error) {
	return r.products, nil
}

// FindByIDs retorna productos por lista de IDs.
func (r *ProductJSONRepository) FindByIDs(ids []string) ([]*domain.Product, error) {
	var result []*domain.Product
	for _, id := range ids {
		product, err := r.FindByID(id)
		if err != nil {
			return nil, err
		}
		result = append(result, product)
	}
	return result, nil
}

// Save agrega un producto al repositorio en memoria.
func (r *ProductJSONRepository) Save(product *domain.Product) (*domain.Product, error) {
	if product == nil {
		return nil, domain.ErrInvalidInput
	}
	r.products = append(r.products, product)
	return product, nil
}
