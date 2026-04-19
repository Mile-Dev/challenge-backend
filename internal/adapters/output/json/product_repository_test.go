package json

import (
	"path/filepath"
	"runtime"
	"testing"

	"project/internal/domain"

	"github.com/stretchr/testify/assert"
)

// getTestDataPath retorna la ruta absoluta al archivo de datos de prueba.
func getTestDataPath() string {
	_, filename, _, _ := runtime.Caller(0)
	dir := filepath.Dir(filename)
	return filepath.Join(dir, "../../../../data/products.json")
}

func TestJSONFindByID_Success(t *testing.T) {
	repo, err := NewProductJSONRepository(getTestDataPath())
	assert.NoError(t, err)

	product, err := repo.FindByID("1")
	assert.NoError(t, err)
	assert.Equal(t, "1", product.ID)
	assert.Equal(t, "Samsung Galaxy S23", product.Name)
}

func TestJSONFindByID_NotFound(t *testing.T) {
	repo, err := NewProductJSONRepository(getTestDataPath())
	assert.NoError(t, err)

	product, err := repo.FindByID("999")
	assert.Nil(t, product)
	assert.ErrorIs(t, err, domain.ErrProductNotFound)
}

func TestJSONFindAll_Success(t *testing.T) {
	repo, err := NewProductJSONRepository(getTestDataPath())
	assert.NoError(t, err)

	products, err := repo.FindAll()
	assert.NoError(t, err)
	assert.Equal(t, 3, len(products))
}

func TestJSONFindByIDs_Success(t *testing.T) {
	repo, err := NewProductJSONRepository(getTestDataPath())
	assert.NoError(t, err)

	products, err := repo.FindByIDs([]string{"1", "2"})
	assert.NoError(t, err)
	assert.Equal(t, 2, len(products))
}

func TestJSONFindByIDs_OneNotFound(t *testing.T) {
	repo, err := NewProductJSONRepository(getTestDataPath())
	assert.NoError(t, err)

	products, err := repo.FindByIDs([]string{"1", "999"})
	assert.Nil(t, products)
	assert.ErrorIs(t, err, domain.ErrProductNotFound)
}

func TestJSONInvalidPath(t *testing.T) {
	repo, err := NewProductJSONRepository("invalid/path.json")
	assert.Nil(t, repo)
	assert.ErrorIs(t, err, domain.ErrInternalServer)
}
