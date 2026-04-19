// @title           Product Comparison API
// @version         1.0
// @description     API RESTful para comparar artículos del catálogo de MercadoLibre
// @host            localhost:8080
// @BasePath        /
package main

import (
	"log"
	"os"

	"project/internal/adapters/input/http"
	jsonrepo "project/internal/adapters/output/json"
	sqliterepo "project/internal/adapters/output/sqlite"
	"project/internal/application"
	"project/internal/ports/output"

	_ "project/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// setupRepository selecciona el repositorio según variable de entorno.
// Por defecto usa JSON, si DB_TYPE=sqlite usa SQLite.
// Esto demuestra el poder de Hexagonal — intercambiamos repositorios
// sin tocar nada de la lógica de negocio.
func setupRepository() output.ProductRepositoryPort {
	dbType := os.Getenv("DB_TYPE")

	switch dbType {
	case "sqlite":
		log.Println("Usando repositorio: SQLite")
		repo, err := sqliterepo.NewProductSQLiteRepository("data/products.db")
		if err != nil {
			log.Fatalf("error iniciando SQLite: %v", err)
		}
		return repo

	default:
		log.Println("Usando repositorio: JSON")
		repo, err := jsonrepo.NewProductJSONRepository("data/products.json")
		if err != nil {
			log.Fatalf("error cargando JSON: %v", err)
		}
		return repo
	}
}

func setupRouter() *gin.Engine {
	// Selecciona repositorio según configuración
	repo := setupRepository()

	// El service y handler no saben qué repositorio están usando
	service := application.NewProductService(repo)
	handler := http.NewProductHandler(service)

	r := gin.Default()
	handler.RegisterRoutes(r)

	// Swagger UI
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return r
}

func main() {
	r := setupRouter()
	r.Run(":8080")
}
