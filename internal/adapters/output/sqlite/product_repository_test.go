package sqlite

import (
	"testing"

	"project/internal/domain"

	"github.com/stretchr/testify/assert"
)

// setupTestRepo crea un repositorio SQLite en memoria para tests.
// :memory: crea una BD temporal que se destruye al terminar el test.
func setupTestRepo(t *testing.T) *ProductSQLiteRepository {
	repo, err := NewProductSQLiteRepository(":memory:")
	assert.NoError(t, err)
	return repo.(*ProductSQLiteRepository)
}

func TestSQLiteFindByID_Success(t *testing.T) {
	repo := setupTestRepo(t)

	product, err := repo.FindByID("1")
	assert.NoError(t, err)
	assert.Equal(t, "1", product.ID)
	assert.Equal(t, "Samsung Galaxy S23", product.Name)
}

func TestSQLiteFindByID_NotFound(t *testing.T) {
	repo := setupTestRepo(t)

	product, err := repo.FindByID("999")
	assert.Nil(t, product)
	assert.ErrorIs(t, err, domain.ErrProductNotFound)
}

func TestSQLiteFindAll_Success(t *testing.T) {
	repo := setupTestRepo(t)

	products, err := repo.FindAll()
	assert.NoError(t, err)
	assert.Greater(t, len(products), 0)
}

func TestSQLiteFindByIDs_Success(t *testing.T) {
	repo := setupTestRepo(t)

	products, err := repo.FindByIDs([]string{"1", "2"})
	assert.NoError(t, err)
	assert.Equal(t, 2, len(products))
}

func TestSQLiteFindByIDs_NotFound(t *testing.T) {
	repo := setupTestRepo(t)

	products, err := repo.FindByIDs([]string{"1", "999"})
	assert.Nil(t, products)
	assert.ErrorIs(t, err, domain.ErrProductNotFound)
}
