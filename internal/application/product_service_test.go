package application

import (
	"errors"
	"testing"

	"project/internal/domain"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockRepository simula el repositorio para tests unitarios.
// Permite probar el servicio sin depender de archivos JSON reales.
type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) FindByID(id string) (*domain.Product, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Product), args.Error(1)
}

func (m *MockRepository) FindAll() ([]*domain.Product, error) {
	args := m.Called()
	return args.Get(0).([]*domain.Product), args.Error(1)
}

func (m *MockRepository) FindByIDs(ids []string) ([]*domain.Product, error) {
	args := m.Called(ids)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Product), args.Error(1)
}

func TestGetProduct_Success(t *testing.T) {
	mockRepo := new(MockRepository)
	service := NewProductService(mockRepo)

	expected := &domain.Product{ID: "1", Name: "Samsung Galaxy S23"}
	mockRepo.On("FindByID", "1").Return(expected, nil)

	result, err := service.GetProduct("1")

	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	mockRepo.AssertExpectations(t)
}

func TestGetProduct_EmptyID(t *testing.T) {
	mockRepo := new(MockRepository)
	service := NewProductService(mockRepo)

	result, err := service.GetProduct("")

	assert.Nil(t, result)
	assert.Equal(t, domain.ErrInvalidInput, err)
}

func TestGetProduct_NotFound(t *testing.T) {
	mockRepo := new(MockRepository)
	service := NewProductService(mockRepo)

	mockRepo.On("FindByID", "999").Return(nil, domain.ErrProductNotFound)

	result, err := service.GetProduct("999")

	assert.Nil(t, result)
	assert.True(t, errors.Is(err, domain.ErrProductNotFound)) // ✅ errors.Is
	mockRepo.AssertExpectations(t)
}

func TestGetAllProducts_Success(t *testing.T) {
	mockRepo := new(MockRepository)
	service := NewProductService(mockRepo)

	expected := []*domain.Product{
		{ID: "1", Name: "Samsung Galaxy S23"},
		{ID: "2", Name: "iPhone 14"},
	}
	mockRepo.On("FindAll").Return(expected, nil)

	result, err := service.GetAllProducts()

	assert.NoError(t, err)
	assert.Equal(t, 2, len(result))
	mockRepo.AssertExpectations(t)
}

func TestCompareProducts_Success(t *testing.T) {
	mockRepo := new(MockRepository)
	service := NewProductService(mockRepo)

	ids := []string{"1", "2"}
	expected := []*domain.Product{
		{ID: "1", Name: "Samsung Galaxy S23"},
		{ID: "2", Name: "iPhone 14"},
	}
	mockRepo.On("FindByIDs", ids).Return(expected, nil)

	result, err := service.CompareProducts(ids)

	assert.NoError(t, err)
	assert.Equal(t, 2, len(result))
	mockRepo.AssertExpectations(t)
}

func TestCompareProducts_LessThanTwoIDs(t *testing.T) {
	mockRepo := new(MockRepository)
	service := NewProductService(mockRepo)

	result, err := service.CompareProducts([]string{"1"})

	assert.Nil(t, result)
	assert.Equal(t, domain.ErrInvalidInput, err)
}
