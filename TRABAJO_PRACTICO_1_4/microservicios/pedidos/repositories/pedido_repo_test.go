package repositories

import (
	"testing"

	"pedidos/models"
)

func TestPedidosMemoria_GuardarAsignaIDCorrelativo(t *testing.T) {
	repo := NuevoPedidosMemoria()

	p1, err := repo.Guardar(models.Pedido{ClienteID: "C-1", ProductoID: "P-1", Cantidad: 1, Estado: "confirmado"})
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	p2, _ := repo.Guardar(models.Pedido{ClienteID: "C-2", ProductoID: "P-2", Cantidad: 3, Estado: "confirmado"})

	if p1.ID != "PED-1" || p2.ID != "PED-2" {
		t.Errorf("IDs = %q, %q; se esperaba PED-1, PED-2", p1.ID, p2.ID)
	}
}

func TestPedidosMemoria_RespetaIDExistente(t *testing.T) {
	repo := NuevoPedidosMemoria()
	p, _ := repo.Guardar(models.Pedido{ID: "PED-XYZ", ClienteID: "C-1", ProductoID: "P-1", Cantidad: 1})
	if p.ID != "PED-XYZ" {
		t.Errorf("ID = %q; se esperaba PED-XYZ", p.ID)
	}
}
