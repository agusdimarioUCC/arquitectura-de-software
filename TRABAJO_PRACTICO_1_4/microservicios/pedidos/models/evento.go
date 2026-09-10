package models

// TipoPedidoConfirmado es el valor del campo "tipo" del evento.
const TipoPedidoConfirmado = "pedido.confirmado"

// EventoPedidoConfirmado es el mensaje que se publica en RabbitMQ
// cuando un pedido queda confirmado.
type EventoPedidoConfirmado struct {
	Tipo       string `json:"tipo"`
	PedidoID   string `json:"pedido_id"`
	ClienteID  string `json:"cliente_id"`
	ProductoID string `json:"producto_id"`
}
