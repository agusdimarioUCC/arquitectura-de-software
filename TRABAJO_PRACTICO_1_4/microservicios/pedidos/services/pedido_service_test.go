package services

import (
	"errors"
	"testing"

	"pedidos/models"
	"pedidos/repositories"
)

// --- dobles de prueba ---

type productosFake struct{ existentes map[string]models.Producto }

func (f productosFake) Listar() ([]models.Producto, error) {
	out := make([]models.Producto, 0, len(f.existentes))
	for _, p := range f.existentes {
		out = append(out, p)
	}
	return out, nil
}
func (f productosFake) BuscarPorID(id string) (models.Producto, error) {
	if p, ok := f.existentes[id]; ok {
		return p, nil
	}
	return models.Producto{}, repositories.ErrProductoNoEncontrado
}

type pedidosFake struct{ guardados []models.Pedido }

func (f *pedidosFake) Guardar(p models.Pedido) (models.Pedido, error) {
	if p.ID == "" {
		p.ID = "PED-TEST"
	}
	f.guardados = append(f.guardados, p)
	return p, nil
}

type publisherFake struct {
	ultimo   models.EventoPedidoConfirmado
	llamado  bool
	devuelve error
}

func (f *publisherFake) PublicarPedidoConfirmado(e models.EventoPedidoConfirmado) error {
	f.llamado = true
	f.ultimo = e
	return f.devuelve
}

func nuevoService(pub *publisherFake) (PedidoService, *pedidosFake) {
	prod := productosFake{existentes: map[string]models.Producto{
		"P-1": {ID: "P-1", Nombre: "Auriculares", Precio: 25000, Stock: 10},
	}}
	ped := &pedidosFake{}
	return NuevoPedidoService(prod, ped, pub), ped
}

// --- tests ---

func TestConfirmar_OkPublicaEvento(t *testing.T) {
	pub := &publisherFake{}
	svc, ped := nuevoService(pub)

	pedido, err := svc.Confirmar("C-1", "P-1", 2)
	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
	if pedido.Estado != "confirmado" {
		t.Errorf("estado = %q; se esperaba confirmado", pedido.Estado)
	}
	if len(ped.guardados) != 1 {
		t.Errorf("pedidos guardados = %d; se esperaba 1", len(ped.guardados))
	}
	if !pub.llamado {
		t.Fatal("no se publicó el evento")
	}
	if pub.ultimo.Tipo != models.TipoPedidoConfirmado ||
		pub.ultimo.PedidoID != pedido.ID ||
		pub.ultimo.ClienteID != "C-1" ||
		pub.ultimo.ProductoID != "P-1" {
		t.Errorf("evento inesperado: %+v", pub.ultimo)
	}
}

func TestConfirmar_ProductoInexistente(t *testing.T) {
	pub := &publisherFake{}
	svc, _ := nuevoService(pub)

	_, err := svc.Confirmar("C-1", "P-99", 1)
	if !errors.Is(err, ErrProductoNoEncontrado) {
		t.Errorf("err = %v; se esperaba ErrProductoNoEncontrado", err)
	}
	if pub.llamado {
		t.Error("no debería publicarse evento si el producto no existe")
	}
}

func TestConfirmar_CantidadInvalida(t *testing.T) {
	pub := &publisherFake{}
	svc, _ := nuevoService(pub)

	_, err := svc.Confirmar("C-1", "P-1", 0)
	if !errors.Is(err, ErrValidacion) {
		t.Errorf("err = %v; se esperaba ErrValidacion", err)
	}
}

func TestConfirmar_CamposVacios(t *testing.T) {
	pub := &publisherFake{}
	svc, _ := nuevoService(pub)

	if _, err := svc.Confirmar("", "P-1", 1); !errors.Is(err, ErrValidacion) {
		t.Errorf("cliente vacío: err = %v", err)
	}
	if _, err := svc.Confirmar("C-1", "  ", 1); !errors.Is(err, ErrValidacion) {
		t.Errorf("producto vacío: err = %v", err)
	}
}

func TestConfirmar_FalloAlPublicar(t *testing.T) {
	pub := &publisherFake{devuelve: errors.New("broker caído")}
	svc, ped := nuevoService(pub)

	pedido, err := svc.Confirmar("C-1", "P-1", 1)
	if !errors.Is(err, ErrPublicacion) {
		t.Errorf("err = %v; se esperaba ErrPublicacion", err)
	}
	if len(ped.guardados) != 1 {
		t.Error("el pedido debería haberse guardado igual antes de publicar")
	}
	if pedido.ID == "" {
		t.Error("se esperaba recuperar el pedido guardado junto al error")
	}
}
