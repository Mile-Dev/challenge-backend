package http

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"project/internal/domain"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockService simula el servicio para tests del handler.
type MockService struct {
	mock.Mock
}

func (m *MockService) GetProduct(id string) (*domain.Product, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Product), args.Error(1)
}

func (m *MockService) GetAllProducts() ([]*domain.Product, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Product), args.Error(1)
}

func (m *MockService) CompareProducts(ids []string) ([]*domain.Product, error) {
	args := m.Called(ids)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Product), args.Error(1)
}

func setupHandlerTest(service *MockService) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	handler := NewProductHandler(service)
	handler.RegisterRoutes(r)
	return r
}

func TestHandlerGetAllProducts_Success(t *testing.T) {
	mockService := new(MockService)
	router := setupHandlerTest(mockService)

	expected := []*domain.Product{
		{ID: "1", Name: "Samsung Galaxy S23"},
		{ID: "2", Name: "iPhone 14"},
	}
	mockService.On("GetAllProducts").Return(expected, nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/products", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var products []domain.Product
	err := json.Unmarshal(w.Body.Bytes(), &products)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(products))
	mockService.AssertExpectations(t)
}

func TestHandlerGetAllProducts_Error(t *testing.T) {
	mockService := new(MockService)
	router := setupHandlerTest(mockService)

	mockService.On("GetAllProducts").Return(nil, domain.ErrInternalServer)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/products", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockService.AssertExpectations(t)
}

func TestHandlerGetProductByID_Success(t *testing.T) {
	mockService := new(MockService)
	router := setupHandlerTest(mockService)

	expected := &domain.Product{ID: "1", Name: "Samsung Galaxy S23"}
	mockService.On("GetProduct", "1").Return(expected, nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/products/1", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestHandlerGetProductByID_NotFound(t *testing.T) {
	mockService := new(MockService)
	router := setupHandlerTest(mockService)

	mockService.On("GetProduct", "999").Return(nil, domain.ErrProductNotFound)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/products/999", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	mockService.AssertExpectations(t)
}

func TestHandlerCompareProducts_Success(t *testing.T) {
	mockService := new(MockService)
	router := setupHandlerTest(mockService)

	ids := []string{"1", "2"}
	expected := []*domain.Product{
		{ID: "1", Name: "Samsung Galaxy S23"},
		{ID: "2", Name: "iPhone 14"},
	}
	mockService.On("CompareProducts", ids).Return(expected, nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/products/compare?ids=1,2", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestHandlerCompareProducts_MissingParam(t *testing.T) {
	mockService := new(MockService)
	router := setupHandlerTest(mockService)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/products/compare", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestHandlerCompareProducts_InvalidInput(t *testing.T) {
	mockService := new(MockService)
	router := setupHandlerTest(mockService)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/products/compare?ids=1", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestHandlerErrorHandler_UnexpectedError(t *testing.T) {
	mockService := new(MockService)
	router := setupHandlerTest(mockService)

	mockService.On("GetProduct", "1").Return(nil, errors.New("unexpected error"))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/products/1", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockService.AssertExpectations(t)
}
