package application

import (
	"fmt"

	"project/internal/domain"
	"project/internal/ports/input"
	"project/internal/ports/output"
)

// ProductService implementa la lógica de negocio para productos.
type ProductService struct {
	repository output.ProductRepositoryPort
}

// NewProductService crea una nueva instancia de ProductService.
func NewProductService(repo output.ProductRepositoryPort) input.ProductServicePort {
	return &ProductService{repository: repo}
}

// GetProduct retorna un producto por su ID.
func (s *ProductService) GetProduct(id string) (*domain.Product, error) {
	if id == "" {
		return nil, domain.ErrInvalidInput
	}

	product, err := s.repository.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("GetProduct id=%s: %w", id, err)
	}
	return product, nil
}

// GetAllProducts retorna todos los productos disponibles.
func (s *ProductService) GetAllProducts() ([]*domain.Product, error) {
	products, err := s.repository.FindAll()
	if err != nil {
		return nil, fmt.Errorf("GetAllProducts: %w", err)
	}
	return products, nil
}

// CompareProducts retorna múltiples productos para comparación.
func (s *ProductService) CompareProducts(ids []string) ([]*domain.Product, error) {
	if len(ids) < 2 {
		return nil, domain.ErrInvalidInput
	}

	products, err := s.repository.FindByIDs(ids)
	if err != nil {
		return nil, fmt.Errorf("CompareProducts ids=%v: %w", ids, err)
	}
	return products, nil
}
