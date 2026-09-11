package messaging

import (
	"encoding/json"
	"testing"

	"pedidos/models"
)

// Verifica en tiempo de compilación que ambas implementaciones cumplen la interfaz.
var (
	_ EventoPublisher = PublisherConsola{}
	_ EventoPublisher = (*RabbitMQPublisher)(nil)
)

func TestPublisherConsola_NoFalla(t *testing.T) {
	err := PublisherConsola{}.PublicarPedidoConfirmado(models.EventoPedidoConfirmado{
		Tipo:       models.TipoPedidoConfirmado,
		PedidoID:   "PED-1",
		ClienteID:  "C-1",
		ProductoID: "P-1",
	})
	if err != nil {
		t.Errorf("err = %v; se esperaba nil", err)
	}
}

func TestNombreCola(t *testing.T) {
	if NombreCola != "pedidos-confirmados" {
		t.Errorf("NombreCola = %q", NombreCola)
	}
}

func TestEvento_SerializaSegunElContrato(t *testing.T) {
	payload, _ := json.Marshal(models.EventoPedidoConfirmado{
		Tipo: models.TipoPedidoConfirmado, PedidoID: "PED-1", ClienteID: "C-1", ProductoID: "P-1",
	})
	esperado := `{"tipo":"pedido.confirmado","pedido_id":"PED-1","cliente_id":"C-1","producto_id":"P-1"}`
	if string(payload) != esperado {
		t.Errorf("payload = %s; se esperaba %s", payload, esperado)
	}
}
