package models

// Producto es un ítem del catálogo.
type Producto struct {
	ID     string  `json:"id"`
	Nombre string  `json:"nombre"`
	Precio float64 `json:"precio"`
	Stock  int     `json:"stock"`
}
