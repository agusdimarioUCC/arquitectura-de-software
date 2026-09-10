package repositories

import (
	"errors"
	"testing"
)

func TestProductosMemoria_ListarDevuelveSeed(t *testing.T) {
	repo := NuevoProductosMemoria()
	lista, err := repo.Listar()
	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
	if len(lista) != 2 {
		t.Fatalf("len(lista) = %d; se esperaba 2", len(lista))
	}
}

func TestProductosMemoria_BuscarPorIDExistente(t *testing.T) {
	repo := NuevoProductosMemoria()
	p, err := repo.BuscarPorID("P-1")
	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
	if p.Nombre != "Auriculares" {
		t.Errorf("nombre = %q; se esperaba %q", p.Nombre, "Auriculares")
	}
}

func TestProductosMemoria_BuscarPorIDInexistente(t *testing.T) {
	repo := NuevoProductosMemoria()
	_, err := repo.BuscarPorID("P-99")
	if !errors.Is(err, ErrProductoNoEncontrado) {
		t.Errorf("err = %v; se esperaba ErrProductoNoEncontrado", err)
	}
}

func TestProductosMemoria_ListarDevuelveCopia(t *testing.T) {
	repo := NuevoProductosMemoria()
	lista, _ := repo.Listar()
	lista[0].Nombre = "MUTADO"
	otra, _ := repo.Listar()
	if otra[0].Nombre == "MUTADO" {
		t.Error("Listar devolvió una referencia mutable al estado interno")
	}
}
