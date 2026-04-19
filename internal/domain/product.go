package domain

// Product representa un artículo del catálogo de MercadoLibre.
// Contiene toda la información necesaria para comparaciones.
type Product struct {
	ID             string            `json:"id"`
	Name           string            `json:"name"`
	Description    string            `json:"description"`
	Price          float64           `json:"price"`
	ImageURL       string            `json:"image_url"`
	Rating         float64           `json:"rating"`
	Specifications map[string]string `json:"specifications"`
}
