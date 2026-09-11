package repositories

import (
	"errors"
	"testing"
)

func TestProductosMemoria_ListarDevuelveSeed(t *testing.T) {
	repositorio := NuevoProductosMemoria()
	lista, err := repositorio.Listar()
	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
	if len(lista) != 2 {
		t.Fatalf("len(lista) = %d; se esperaba 2", len(lista))
	}
}

func TestProductosMemoria_BuscarPorIDExistente(t *testing.T) {
	repositorio := NuevoProductosMemoria()
	producto, err := repositorio.BuscarPorID("P-1")
	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
	if producto.Nombre != "Auriculares" {
		t.Errorf("nombre = %q; se esperaba %q", producto.Nombre, "Auriculares")
	}
}

func TestProductosMemoria_BuscarPorIDInexistente(t *testing.T) {
	repositorio := NuevoProductosMemoria()
	_, err := repositorio.BuscarPorID("P-99")
	if !errors.Is(err, ErrProductoNoEncontrado) {
		t.Errorf("err = %v; se esperaba ErrProductoNoEncontrado", err)
	}
}

func TestProductosMemoria_ListarDevuelveCopia(t *testing.T) {
	repositorio := NuevoProductosMemoria()
	lista, _ := repositorio.Listar()
	lista[0].Nombre = "MUTADO"
	otraLista, _ := repositorio.Listar()
	if otraLista[0].Nombre == "MUTADO" {
		t.Error("Listar devolvió una referencia mutable al estado interno")
	}
}
