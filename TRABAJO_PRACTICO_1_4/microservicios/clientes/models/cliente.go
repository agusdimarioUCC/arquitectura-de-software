package models

// Cliente representa a un cliente del e-commerce.
type Cliente struct {
	ID     string `json:"id"`
	Nombre string `json:"nombre"`
	Email  string `json:"email"`
}
