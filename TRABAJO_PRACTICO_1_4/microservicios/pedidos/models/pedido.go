package models

// Pedido representa una compra confirmada.
type Pedido struct {
	ID         string `json:"pedido_id"`
	ClienteID  string `json:"cliente_id"`
	ProductoID string `json:"producto_id"`
	Cantidad   int    `json:"cantidad"`
	Estado     string `json:"estado"`
}
