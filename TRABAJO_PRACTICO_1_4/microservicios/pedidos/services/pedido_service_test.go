package services

import (
	"errors"
	"testing"

	"pedidos/models"
	"pedidos/repositories"
)

// --- dobles de prueba ---

type productosSimulados struct{ existentes map[string]models.Producto }

func (simulado productosSimulados) Listar() ([]models.Producto, error) {
	lista := make([]models.Producto, 0, len(simulado.existentes))
	for _, producto := range simulado.existentes {
		lista = append(lista, producto)
	}
	return lista, nil
}
func (simulado productosSimulados) BuscarPorID(identificador string) (models.Producto, error) {
	if producto, existe := simulado.existentes[identificador]; existe {
		return producto, nil
	}
	return models.Producto{}, repositories.ErrProductoNoEncontrado
}

type pedidosSimulados struct{ guardados []models.Pedido }

func (simulado *pedidosSimulados) Guardar(pedido models.Pedido) (models.Pedido, error) {
	if pedido.ID == "" {
		pedido.ID = "PED-TEST"
	}
	simulado.guardados = append(simulado.guardados, pedido)
	return pedido, nil
}

type publicadorSimulado struct {
	ultimoEvento models.EventoPedidoConfirmado
	llamado      bool
	devuelve     error
}

func (simulado *publicadorSimulado) PublicarPedidoConfirmado(evento models.EventoPedidoConfirmado) error {
	simulado.llamado = true
	simulado.ultimoEvento = evento
	return simulado.devuelve
}

func nuevoServicio(publicador *publicadorSimulado) (PedidoService, *pedidosSimulados) {
	productos := productosSimulados{existentes: map[string]models.Producto{
		"P-1": {ID: "P-1", Nombre: "Auriculares", Precio: 25000, Stock: 10},
	}}
	pedidos := &pedidosSimulados{}
	return NuevoPedidoService(productos, pedidos, publicador), pedidos
}

// --- tests ---

func TestConfirmar_OkPublicaEvento(t *testing.T) {
	publicador := &publicadorSimulado{}
	servicio, pedidosGuardados := nuevoServicio(publicador)

	pedido, err := servicio.Confirmar("C-1", "P-1", 2)
	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
	if pedido.Estado != "confirmado" {
		t.Errorf("estado = %q; se esperaba confirmado", pedido.Estado)
	}
	if len(pedidosGuardados.guardados) != 1 {
		t.Errorf("pedidos guardados = %d; se esperaba 1", len(pedidosGuardados.guardados))
	}
	if !publicador.llamado {
		t.Fatal("no se publicó el evento")
	}
	if publicador.ultimoEvento.Tipo != models.TipoPedidoConfirmado ||
		publicador.ultimoEvento.PedidoID != pedido.ID ||
		publicador.ultimoEvento.ClienteID != "C-1" ||
		publicador.ultimoEvento.ProductoID != "P-1" {
		t.Errorf("evento inesperado: %+v", publicador.ultimoEvento)
	}
}

func TestConfirmar_ProductoInexistente(t *testing.T) {
	publicador := &publicadorSimulado{}
	servicio, _ := nuevoServicio(publicador)

	_, err := servicio.Confirmar("C-1", "P-99", 1)
	if !errors.Is(err, ErrProductoNoEncontrado) {
		t.Errorf("err = %v; se esperaba ErrProductoNoEncontrado", err)
	}
	if publicador.llamado {
		t.Error("no debería publicarse evento si el producto no existe")
	}
}

func TestConfirmar_CantidadInvalida(t *testing.T) {
	publicador := &publicadorSimulado{}
	servicio, _ := nuevoServicio(publicador)

	_, err := servicio.Confirmar("C-1", "P-1", 0)
	if !errors.Is(err, ErrValidacion) {
		t.Errorf("err = %v; se esperaba ErrValidacion", err)
	}
}

func TestConfirmar_CamposVacios(t *testing.T) {
	publicador := &publicadorSimulado{}
	servicio, _ := nuevoServicio(publicador)

	if _, err := servicio.Confirmar("", "P-1", 1); !errors.Is(err, ErrValidacion) {
		t.Errorf("cliente vacío: err = %v", err)
	}
	if _, err := servicio.Confirmar("C-1", "  ", 1); !errors.Is(err, ErrValidacion) {
		t.Errorf("producto vacío: err = %v", err)
	}
}

func TestConfirmar_FalloAlPublicar(t *testing.T) {
	publicador := &publicadorSimulado{devuelve: errors.New("broker caído")}
	servicio, pedidosGuardados := nuevoServicio(publicador)

	pedido, err := servicio.Confirmar("C-1", "P-1", 1)
	if !errors.Is(err, ErrPublicacion) {
		t.Errorf("err = %v; se esperaba ErrPublicacion", err)
	}
	if len(pedidosGuardados.guardados) != 1 {
		t.Error("el pedido debería haberse guardado igual antes de publicar")
	}
	if pedido.ID == "" {
		t.Error("se esperaba recuperar el pedido guardado junto al error")
	}
}
