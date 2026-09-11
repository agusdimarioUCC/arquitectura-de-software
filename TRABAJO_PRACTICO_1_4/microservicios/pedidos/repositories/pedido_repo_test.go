package repositories

import (
	"testing"

	"pedidos/models"
)

func TestPedidosMemoria_GuardarAsignaIDCorrelativo(t *testing.T) {
	repositorio := NuevoPedidosMemoria()

	primerPedido, err := repositorio.Guardar(models.Pedido{ClienteID: "C-1", ProductoID: "P-1", Cantidad: 1, Estado: "confirmado"})
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	segundoPedido, _ := repositorio.Guardar(models.Pedido{ClienteID: "C-2", ProductoID: "P-2", Cantidad: 3, Estado: "confirmado"})

	if primerPedido.ID != "PED-1" || segundoPedido.ID != "PED-2" {
		t.Errorf("IDs = %q, %q; se esperaba PED-1, PED-2", primerPedido.ID, segundoPedido.ID)
	}
}

func TestPedidosMemoria_RespetaIDExistente(t *testing.T) {
	repositorio := NuevoPedidosMemoria()
	pedido, _ := repositorio.Guardar(models.Pedido{ID: "PED-XYZ", ClienteID: "C-1", ProductoID: "P-1", Cantidad: 1})
	if pedido.ID != "PED-XYZ" {
		t.Errorf("ID = %q; se esperaba PED-XYZ", pedido.ID)
	}
}
