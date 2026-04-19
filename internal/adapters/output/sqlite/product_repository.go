package sqlite

import (
	"database/sql"
	"encoding/json"

	"project/internal/domain"
	"project/internal/ports/output"

	_ "modernc.org/sqlite"
)

// ProductSQLiteRepository implementa ProductRepositoryPort usando SQLite.
// Demuestra que el puerto de salida es intercambiable sin tocar el dominio.
type ProductSQLiteRepository struct {
	db *sql.DB
}

// NewProductSQLiteRepository crea una nueva instancia y prepara la base de datos.
func NewProductSQLiteRepository(dbPath string) (output.ProductRepositoryPort, error) {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, domain.ErrInternalServer
	}

	repo := &ProductSQLiteRepository{db: db}

	if err := repo.createTable(); err != nil {
		return nil, domain.ErrInternalServer
	}

	if err := repo.seedData(); err != nil {
		return nil, domain.ErrInternalServer
	}

	return repo, nil
}

// createTable crea la tabla de productos si no existe.
func (r *ProductSQLiteRepository) createTable() error {
	query := `
	CREATE TABLE IF NOT EXISTS products (
		id TEXT PRIMARY KEY,
		name TEXT NOT NULL,
		description TEXT,
		price REAL,
		image_url TEXT,
		rating REAL,
		specifications TEXT
	)`
	_, err := r.db.Exec(query)
	return err
}

// seedData inserta datos iniciales si la tabla está vacía.
func (r *ProductSQLiteRepository) seedData() error {
	var count int
	r.db.QueryRow("SELECT COUNT(*) FROM products").Scan(&count)
	if count > 0 {
		return nil
	}

	products := []domain.Product{
		{
			ID:          "1",
			Name:        "Samsung Galaxy S23",
			Description: "Smartphone flagship con cámara de 50MP",
			Price:       2999999.00,
			ImageURL:    "https://example.com/samsung-s23.jpg",
			Rating:      4.5,
			Specifications: map[string]string{
				"ram":     "8GB",
				"storage": "128GB",
				"battery": "3900mAh",
				"screen":  "6.1 pulgadas",
			},
		},
		{
			ID:          "2",
			Name:        "iPhone 14",
			Description: "Smartphone Apple con chip A15 Bionic",
			Price:       3499999.00,
			ImageURL:    "https://example.com/iphone14.jpg",
			Rating:      4.7,
			Specifications: map[string]string{
				"ram":     "6GB",
				"storage": "128GB",
				"battery": "3279mAh",
				"screen":  "6.1 pulgadas",
			},
		},
		{
			ID:          "3",
			Name:        "Xiaomi Redmi Note 12",
			Description: "Smartphone gama media con pantalla AMOLED",
			Price:       999999.00,
			ImageURL:    "https://example.com/redmi-note12.jpg",
			Rating:      4.2,
			Specifications: map[string]string{
				"ram":     "6GB",
				"storage": "128GB",
				"battery": "5000mAh",
				"screen":  "6.67 pulgadas",
			},
		},
		{
			ID:          "4",
			Name:        "iPhone 15 Pro",
			Description: "Smartphone Apple con chip A17 Pro y titanio",
			Price:       5999999.00,
			ImageURL:    "https://example.com/iphone15pro.jpg",
			Rating:      4.9,
			Specifications: map[string]string{
				"ram":     "8GB",
				"storage": "256GB",
				"battery": "3274mAh",
				"screen":  "6.1 pulgadas",
			},
		},
		{
			ID:          "5",
			Name:        "iPhone 13",
			Description: "Smartphone Apple con chip A15 y modo cinematico",
			Price:       2499999.00,
			ImageURL:    "https://example.com/iphone13.jpg",
			Rating:      4.6,
			Specifications: map[string]string{
				"ram":     "4GB",
				"storage": "128GB",
				"battery": "3227mAh",
				"screen":  "6.1 pulgadas",
			},
		},
		{
			ID:          "6",
			Name:        "iPhone SE 2022",
			Description: "Smartphone Apple compacto con chip A15 Bionic",
			Price:       1799999.00,
			ImageURL:    "https://example.com/iphonese.jpg",
			Rating:      4.3,
			Specifications: map[string]string{
				"ram":     "4GB",
				"storage": "64GB",
				"battery": "2018mAh",
				"screen":  "4.7 pulgadas",
			},
		},
	}

	for _, p := range products {
		specs, err := json.Marshal(p.Specifications)
		if err != nil {
			return err
		}
		_, err = r.db.Exec(
			`INSERT INTO products 
			(id, name, description, price, image_url, rating, specifications) 
			VALUES (?, ?, ?, ?, ?, ?, ?)`,
			p.ID, p.Name, p.Description,
			p.Price, p.ImageURL, p.Rating,
			string(specs),
		)
		if err != nil {
			return err
		}
	}
	return nil
}

// FindByID busca un producto por su ID.
func (r *ProductSQLiteRepository) FindByID(id string) (*domain.Product, error) {
	query := `SELECT id, name, description, price, image_url, rating, specifications 
			  FROM products WHERE id = ?`
	row := r.db.QueryRow(query, id)
	product, err := r.scanProduct(row)
	if err != nil {
		return nil, err
	}
	return product, nil
}

// FindAll retorna todos los productos.
func (r *ProductSQLiteRepository) FindAll() ([]*domain.Product, error) {
	query := `SELECT id, name, description, price, image_url, rating, specifications 
			  FROM products`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, domain.ErrInternalServer
	}
	defer rows.Close()

	var products []*domain.Product
	for rows.Next() {
		product, err := r.scanProduct(rows)
		if err != nil {
			return nil, err
		}
		products = append(products, product)
	}
	return products, nil
}

// FindByIDs retorna productos por lista de IDs.
func (r *ProductSQLiteRepository) FindByIDs(ids []string) ([]*domain.Product, error) {
	var products []*domain.Product
	for _, id := range ids {
		product, err := r.FindByID(id)
		if err != nil {
			return nil, err
		}
		products = append(products, product)
	}
	return products, nil
}

// scanProduct convierte una fila de SQLite en un Product.
func (r *ProductSQLiteRepository) scanProduct(scanner interface {
	Scan(dest ...interface{}) error
}) (*domain.Product, error) {
	var p domain.Product
	var specsJSON string

	err := scanner.Scan(
		&p.ID,
		&p.Name,
		&p.Description,
		&p.Price,
		&p.ImageURL,
		&p.Rating,
		&specsJSON,
	)
	if err == sql.ErrNoRows {
		return nil, domain.ErrProductNotFound
	}
	if err != nil {
		return nil, domain.ErrInternalServer
	}

	if err := json.Unmarshal([]byte(specsJSON), &p.Specifications); err != nil {
		return nil, domain.ErrInternalServer
	}

	return &p, nil
}
